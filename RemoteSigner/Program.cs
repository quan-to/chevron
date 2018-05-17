using RemoteSigner.Database;
using RemoteSigner.Log;

namespace RemoteSigner {

    class MainClass {

        public static void Main(string[] args) {
            Logger.GlobalEnableDebug = true;
            var dm = DatabaseManager.GlobalDm.GetConnection();
            Logger.Log("Application", $"Database Hostname: {dm.Hostname}");
            RancherThread rt = new RancherThread();
            RancherManager.Init();
            if (RancherManager.InRancher) {
                Logger.Log("Application", "Running in rancher. Starting Rancher Sentinel.");
                rt.Start();
            }
            Http httpServer = new Http(Configuration.HttpPort);
            httpServer.StartSync();
        }
    }
}
