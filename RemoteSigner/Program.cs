using RemoteSigner.Log;

namespace RemoteSigner {

    class MainClass {

        public static void Main(string[] args) {
            Logger.GlobalEnableDebug = true;
            Http httpServer = new Http();
            httpServer.StartSync();
        }
    }
}
