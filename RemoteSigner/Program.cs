using RemoteSigner.Log;

namespace RemoteSigner {

    class MainClass {

        public static void Main(string[] args) {
            Logger.GlobalEnableDebug = true;
            RancherThread rt = new RancherThread();
            if (RancherManager.InRancher) {
                Logger.Log("Application", "Running in rancher. Starting Rancher Sentinel.");
                rt.Start();
            }
            Http httpServer = new Http(Configuration.HttpPort);
            httpServer.StartSync();
        }
    }
}
