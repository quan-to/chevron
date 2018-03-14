using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using RemoteSigner.Log;
using RemoteSigner.Models;

namespace RemoteSigner {
    public static class RancherManager {

        const string RancherMetadata = "http://rancher-metadata/2015-12-19";
        const string RancherManagerLog = "RancherManager";

        public static bool InRancher { get; private set; }
        public static string UUID { get; private set; }

        static RancherManager() {

        }

        public static void Init() {
            Logger.Log(RancherManagerLog, "Checking if running in rancher");
            var check = CheckInRancher();
            check.Wait();
            InRancher = check.Result;

            if (InRancher) {
                var uuidTask = MyUUID();
                uuidTask.Wait();
                UUID = uuidTask.Result;
            }
        }

        static async Task<bool> CheckInRancher() {
            try {
                await MyUUID();
                return true;
            } catch (Exception e) {
                Logger.Warn(RancherManagerLog, $"Probably not in Rancher Mode: {e.Message}");
                Logger.Warn(RancherManagerLog, e.StackTrace);
                return false;
            }
        }

        static async Task<string> MyUUID() {
            return await Tools.Get($"{RancherMetadata}/self/container/uuid/");
        }

        public static async Task<List<RancherNode>> GetServiceNodes() {
            var nodesString = await Tools.Get($"{RancherMetadata}/self/service/containers");
            var nodes = nodesString.Split('\n').Where(x => !string.IsNullOrEmpty(x)).ToList();
            var results = await Task.WhenAll(nodes.Select(async (nodeS) => {
                Logger.Log(RancherManagerLog, $"Checking node {nodeS}");
                var nodeData = nodeS.Split('=');
                var uuid = await Tools.Get($"{RancherMetadata}/self/service/containers/{nodeData[0]}/uuid/");
                var ipaddr = await Tools.Get($"{RancherMetadata}/self/service/containers/{nodeData[0]}/primary_ip/");
                return new RancherNode {
                    UUID = uuid,
                    IsSelf = uuid == UUID,
                    Name = nodeData[1],
                    ID = int.Parse(nodeData[0]),
                    IPAddress = ipaddr,
                };
            }));
            return results.ToList();
        }
    }
}
