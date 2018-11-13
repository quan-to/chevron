using RemoteSigner.Database;
using RemoteSigner.Log;

namespace RemoteSigner {
    static class MainClass {
        public static void Main(string[] args) {
            Logger.GlobalEnableDebug = true;
            if (Configuration.EnableRethinkSKS) {
                var dm = DatabaseManager.GlobalDm.GetConnection();
                Logger.Log("Application", $"Database Hostname: {dm.Hostname}");
            }
            
            RancherManager.Init();
            Kubernetes.Init();

            if (RancherManager.InRancher) {
                Logger.Log("Application", "Running in rancher. Starting Rancher Sentinel.");
                var rt = new RancherThread();
                rt.Start();
            }

            if (Kubernetes.InKubernetes) {
                Logger.Log("Application", "Running in kubernetes. Starting Kubernetes Sentinel.");
                var rt = new KubernetesThread();
                rt.Start();
            }
            var httpServer = new Http(Configuration.HttpPort);
            httpServer.StartSync();
        }
    }
}
