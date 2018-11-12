using System;
using System.Collections.Generic;
using System.IO;
using System.Net;
using System.Threading.Tasks;
using Newtonsoft.Json;
using RemoteSigner.Log;
using RemoteSigner.Models.Kubernetes;

namespace RemoteSigner {
    public class Kubernetes {
        private const string ServiceAccountPath = "/run/secrets/kubernetes.io/serviceaccount";
        private const string KubernetesLog = "Kubernetes";

        private static JsonSerializerSettings jsonSettings;
        
        public static bool InKubernetes { get; private set; }
        public static string Namespace { get; private set; }
        public static Pod Me { get; private set; }
        public static string Hostname { get; private set; }

        private static string ServiceURL => $"https://kubernetes.default.svc/api/v1/namespaces/{Namespace}/services";
        private static string PodURL => $"https://kubernetes.default.svc/api/v1/namespaces/{Namespace}/pods";

        private static string KubeToken { get; set; }

        public static void Init() {
            Logger.Log(KubernetesLog, "Checking if running in kubernetes");
            InKubernetes = CheckInKubernetes();

            if (!InKubernetes) return;

            jsonSettings = new JsonSerializerSettings {
                MissingMemberHandling = MissingMemberHandling.Ignore,
            };


            Hostname = Dns.GetHostName();
            Namespace = MyNamespace();
            KubeToken = LoadKubeToken();
            var mTask = MySelf();
            mTask.Wait();
            Me = mTask.Result;
        }

        private static bool CheckInKubernetes() {
            try {
                return File.Exists($"{ServiceAccountPath}/token");
            } catch (Exception e) {
                Logger.Warn(KubernetesLog, $"Probably not in Kubernetes Mode: {e.Message}");
                return false;
            }
        }

        private static string MyNamespace() {
            return File.ReadAllText($"{ServiceAccountPath}/namespace");
        }

        private static string LoadKubeToken() {
            return File.ReadAllText($"{ServiceAccountPath}/token");
        }

        private static async Task<Pod> MySelf() {
            var data = await Tools.Get($"{PodURL}/{Hostname}", KubeToken);
            return JsonConvert.DeserializeObject<Pod>(data, jsonSettings);
        }

        public static async Task<Service> Services() {
            var data = await Tools.Get(ServiceURL, KubeToken);
            return JsonConvert.DeserializeObject<Service>(data, jsonSettings);
        }

        public static async Task<Pods> Pods() {
            var data = await Tools.Get(PodURL, KubeToken);
            return JsonConvert.DeserializeObject<Pods>(data, jsonSettings);
        }
    }
}
