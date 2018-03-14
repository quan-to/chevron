using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Newtonsoft.Json;
using RemoteSigner.Log;

namespace RemoteSigner {
    public class RancherThread {

        const string RancherThreadLog = "RancherSentinel";
        const int SLEEP_INTERVAL = 1 * 60 * 1000;

        bool running;
        Thread thread;

        public void Start() {
            if (!running) {
                Logger.Log(RancherThreadLog, "Starting Rancher Thread");
                running = true;
                thread = new Thread(ThreadLoop) {
                    IsBackground = true
                };
                thread.Start();
            } else {
                Logger.Log(RancherThreadLog, "Rancher Thread already running...");
            }
        }

        public void Stop() {
            if (running) {
                Logger.Log(RancherThreadLog, "Stopping Rancher Thread");
                running = false;
            } else {
                Logger.Log(RancherThreadLog, "Rancher Thread is already stopped");
            }
        }

        void ThreadLoop() {
            int randomWaitTime = (new Random().Next(5)) * 1000 + 1 * 1000;
            Logger.Log(RancherThreadLog, "To avoid concurrency on cluster starting we're waiting 5 seconds plus some random time");
            Logger.Log(RancherThreadLog, $"The exact time is {randomWaitTime} ms");
            Thread.Sleep(randomWaitTime);
            while (running) {
                Logger.Log(RancherThreadLog, "Checking for other RemoteSigner nodes...");
                var task = ThreadFunc();
                task.Wait();
                Logger.Log(RancherThreadLog, $"Sleeping for {SLEEP_INTERVAL} ms.");
                Thread.Sleep(SLEEP_INTERVAL);
            }
        }

        async Task ThreadFunc() {
            try {
                var nodes = await RancherManager.GetServiceNodes();
                Logger.Log(RancherThreadLog, $"There are {nodes.Count} nodes available (including me). Fetching encrypted passwords.");
                foreach (var node in nodes) {
                    try {
                        if (!node.IsSelf) {
                            var getUrl = $"http://{node.IPAddress}:5100/remoteSigner/__internal/__getUnlockPasswords";
                            var postUrl = $"http://localhost:5100/remoteSigner/__internal/__postEncryptedPasswords";
                            var jsonData = await Tools.Get(getUrl);
                            var keys = JsonConvert.DeserializeObject<Dictionary<string, string>>(jsonData);
                            Logger.Log(RancherThreadLog, $"Got {keys.Keys.Count} passwords for node {node.IPAddress}. Adding to local password holder.");
                            await Tools.Post(postUrl, jsonData);
                        }
                    } catch (Exception e) {
                        Logger.Error(RancherThreadLog, $"Error checking node passwords: {e.Message}");
                    }
                }
                Logger.Log(RancherThreadLog, $"Passwords added. Triggering local unlock.");
                await Tools.Get("http://localhost:5100/remoteSigner/__internal/__triggerKeyUnlock");
            } catch (Exception e) {
                Logger.Error(RancherThreadLog, $"There was an error checking for other nodes: {e.Message}");
            }
        }
    }
}
