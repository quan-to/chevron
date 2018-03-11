using RemoteSigner.Log;

namespace RemoteSigner {

    class MainClass {

        public static void Main(string[] args) {
            Logger.GlobalEnableDebug = true;
            if (RancherManager.InRancher) {
                Logger.Log("Application", "Running in rancher. Starting Rancher Sentinel.");
            }
            Http httpServer = new Http(Configuration.HttpPort);
            httpServer.StartSync();
        }
    }
}
