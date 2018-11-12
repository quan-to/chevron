using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Newtonsoft.Json;
using RemoteSigner.Log;
using RemoteSigner.Models.Kubernetes;

namespace RemoteSigner {
    public class KubernetesThread {
        private const string KubernetesThreadLog = "KubernetesSentinel";
        private const int SleepInterval = 1 * 60 * 1000;

        bool running;
        Thread thread;

        public void Start() {
            if (!running) {
                Logger.Log(KubernetesThreadLog, "Starting Kubernetes Thread");
                running = true;
                thread = new Thread(ThreadLoop) {
                    IsBackground = true
                };
                thread.Start();
            } else {
                Logger.Log(KubernetesThreadLog, "Kubernetes Thread already running...");
            }
        }

        public void Stop() {
            if (running) {
                Logger.Log(KubernetesThreadLog, "Stopping Kubernetes Thread");
                running = false;
            } else {
                Logger.Log(KubernetesThreadLog, "Kubernetes Thread is already stopped");
            }
        }

        private void ThreadLoop() {
            var randomWaitTime = (new Random().Next(5)) * 1000 + 1 * 1000;
            Logger.Log(KubernetesThreadLog, $"Kubernetes Namespace {Kubernetes.Namespace}");
            Logger.Log(KubernetesThreadLog, $"Pod Hostname {Kubernetes.Hostname}");
            Logger.Log(KubernetesThreadLog, $"Pod ID {Kubernetes.Me.Metadata.UID}");
            Logger.Log(KubernetesThreadLog, "To avoid concurrency on cluster starting we're waiting 1 second plus some random time");
            Logger.Log(KubernetesThreadLog, $"The exact time is {randomWaitTime} ms");
            Thread.Sleep(randomWaitTime);
            while (running) {
                Logger.Log(KubernetesThreadLog, "Checking for other RemoteSigner pods...");
                var task = ThreadFunc();
                task.Wait();
                Logger.Log(KubernetesThreadLog, $"Sleeping for {SleepInterval} ms.");
                Thread.Sleep(SleepInterval);
            }
        }

        private static async Task ThreadFunc() {
            try {
                var pods = await Kubernetes.Pods();
                var myId = Kubernetes.Me.Metadata.UID;
                Logger.Log(KubernetesThreadLog, $"There are {pods.Items.Count} pods available (including me). Fetching encrypted passwords.");
                foreach (var pod in pods.Items) {
                    try {
                        if (pod.Metadata.UID == myId) continue;
                        if (pod.Status.Phase != Phase.Running) continue;
                        var getUrl = $"http://{pod.Status.PodIP}:{Configuration.HttpPort}/remoteSigner/__internal/__getUnlockPasswords";
                        var postUrl = $"http://localhost:{Configuration.HttpPort}/remoteSigner/__internal/__postEncryptedPasswords";
                        var jsonData = await Tools.Get(getUrl);
                        var keys = JsonConvert.DeserializeObject<Dictionary<string, string>>(jsonData);
                        Logger.Log(KubernetesThreadLog, $"Got {keys.Keys.Count} passwords for node {pod.Status.PodIP}. Adding to local password holder.");
                        await Tools.Post(postUrl, jsonData);
                    } catch (Exception e) {
                        Logger.Error(KubernetesThreadLog, $"Error checking node passwords: {e.Message}");
                    }
                }
                Logger.Log(KubernetesThreadLog, $"Passwords added. Triggering local unlock.");
                await Tools.Get($"http://localhost:{Configuration.HttpPort}/remoteSigner/__internal/__triggerKeyUnlock");
            } catch (Exception e) {
                Logger.Error(KubernetesThreadLog, $"There was an error checking for other nodes: {e.Message}");
            }
        }
    }
}
