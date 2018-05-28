using System;
using System.Linq;
using System.Net;
using System.Reflection;
using System.Text;
using System.Threading;
using RemoteSigner.AppServer;
using RemoteSigner.Log;
using RemoteSigner.Models;

namespace RemoteSigner {
    public class Http {
        public int Port { get; private set; }

        readonly HttpListener listener = new HttpListener();
        Thread listenerThread;
        bool running;
        readonly RestProcessor restProcessor;

        public Http(int port = 5100) {
            Port = port;
            listener.Prefixes.Add($"http://*:{port}/");
            listenerThread = null;
            running = false;
            restProcessor = new RestProcessor();
            LoadEndpoints();
        }

        public void Start() {
            if (listenerThread != null) {
                Logger.Log("HTTP Server", "Starting HTTP Listener");
                listener.Start();
                listenerThread = new Thread(new ThreadStart(ListenerProcessor)) {
                    IsBackground = true
                };
                running = true;
                listenerThread.Start();
            } else {
                Logger.Error("HTTP Server", "HTTP Thread already started!");
            }
        }

        public void StartSync() {
            listener.Start();
            running = true;
            ListenerProcessor();
        }

        public void Stop() {
            running = false;
            if (listenerThread != null) {
                listenerThread.Join();
                listenerThread = null;
            }
        }

        void LoadEndpoints() {
            Assembly a = Assembly.GetExecutingAssembly();
            string[] namespaces = a.GetTypes().Select(x => x.Namespace).Distinct().ToArray();
            foreach (string n in namespaces) {
                if (n.StartsWith("RemoteSigner.HttpData", StringComparison.InvariantCultureIgnoreCase)) {
                    Logger.Log("HTTP Server", $"Loading REST calls for {n}");
                    restProcessor.Init(a, n);
                }
            }
        }

        void ListenerProcessor() {
            while (running) {
                try {
                    var context = listener.GetContext();
                    ThreadPool.QueueUserWorkItem(o => HandleRequest(context));
                } catch (Exception e) {
                    Logger.Error($"Error handling HTTP Request: {e}");
                }
            }
            listener.Stop();
        }

        void HandleRequest(object state) {
            var ctx = state as HttpListenerContext;
            try {
                // Logger.Debug("HTTP Server", $"{ctx.Request.HttpMethod} - {ctx.Request.RawUrl}");
                RestResult ret = ProcessHttpCalls(ctx.Request);
                ctx.Response.ContentType = ret.ContentType;
                ctx.Response.StatusCode = (int)ret.StatusCode;
                ctx.Response.ContentLength64 = ret.Result.Length;
                ctx.Response.OutputStream.Write(ret.Result, 0, ret.Result.Length);
            } catch (Exception e) {
                Logger.Error("HTTP Server", $"Error processing HTTP Request: {e}");
            } finally {
                ctx.Response.OutputStream.Close();
            }
        }

        RestResult ProcessHttpCalls(HttpListenerRequest request) {
            string[] ePath = request.Url.AbsolutePath.Split(new char[] { '/' }, 2, StringSplitOptions.RemoveEmptyEntries);
            RestRequest req = new RestRequest(request);
            if (ePath.Length == 0) {
                return new RestResult(new ErrorObject {
                    ErrorCode = ErrorCodes.NotFound,
                    Message = "No application specified",
                    ErrorField = "url"
                }.ToJSON(), MimeTypes.JSON, HttpStatusCode.NotFound);
            }
            string path = ePath.Length > 1 ? "/" + ePath[1] : "/";
            string method = request.HttpMethod;
            string app = ePath[0];

            if (app != "remoteSigner") {
                /*return new RestResult(new ErrorObject {
                    ErrorCode = ErrorCodes.NotFound,
                    Message = $"Application {app} not found",
                    ErrorField = "url"
                }.ToJSON(), MimeTypes.JSON, HttpStatusCode.NotFound);*/
                // Bypass
                path = "/" + ePath[0] + path;
            }

            Logger.Debug("HTTP Server", $"Processing HTTP Call for App {app}: {method} {path}");

            if (restProcessor.ContainsEndPoint(path, method)) {
                try {
                    return restProcessor.CallEndPoint(path, method, req);
                } catch (Exception e) {
                    string exceptionName = e.InnerException.GetType().Name;
                    string baseName = e.InnerException.GetType().BaseType.Name;
                    IRestExceptionHandler handler = restProcessor.GetExceptionHandler(exceptionName) ?? restProcessor.GetExceptionHandler(baseName);
                    if (handler != null) {
                        return handler.HandleException(e.InnerException);
                    }
                    RestResult result = new RestResult();
                    string exceptionMessage;
                    if (e.InnerException != null) { // In the rest exceptions the real exception will be at InnerException.
                        Logger.Debug("HTTP Server", $"Exception when calling application {app} in endpoint {method} {path} \r\n {e.InnerException}");
                        exceptionMessage = e.InnerException.ToString();
                    } else { // But if we got a internal exception at AppServer, it will be in the root.
                        Logger.Debug("HTTP Server", $"Exception when calling application {app} in endpoint {method} {path} \r\n {e}");
                        exceptionMessage = e.ToString();
                    }
                    result.StatusCode = HttpStatusCode.InternalServerError;
                    result.ContentType = "text/plain";
                    result.Result = Encoding.UTF8.GetBytes(exceptionMessage);
                    return result;
                }
            }
            return new RestResult(new ErrorObject {
                ErrorCode = ErrorCodes.NotFound,
                Message = "Endpoint not found",
                ErrorField = "url"
            }.ToJSON(), MimeTypes.JSON, HttpStatusCode.NotFound);
        }
    }
}
