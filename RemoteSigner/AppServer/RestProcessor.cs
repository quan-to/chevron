using System;
using System.Collections.Generic;
using System.Net;
using System.Reflection;
using Newtonsoft.Json;
using RemoteSigner.Exceptions;
using RemoteSigner.Log;
using RemoteSigner.Models;
using RemoteSigner.Models.Attributes;

namespace RemoteSigner.AppServer {
    public class RestProcessor {
        readonly Type[] RestTypes = { typeof(GET), typeof(POST), typeof(PUT), typeof(DELETE) };
        readonly Dictionary<string, Dictionary<string, RestCall>> endpoints;
        Dictionary<string, RestProxy> proxies;
        Dictionary<string, IRestExceptionHandler> exceptionHandlers;
        Dictionary<string, Object> injectables;

        public RestProcessor() {
            endpoints = new Dictionary<string, Dictionary<string, RestCall>>();
            proxies = new Dictionary<string, RestProxy>();
            exceptionHandlers = new Dictionary<string, IRestExceptionHandler>();
            injectables = new Dictionary<string, object>();
        }

        public void Init(Assembly runningAssembly, string modulesAssembly) {
            Type[] typelist = Tools.GetTypesInNamespace(runningAssembly, modulesAssembly);
            for (int i = 0; i < typelist.Length; i++) {
                Type tClass = typelist[i];

                // Search for REST Attribute
                Attribute t = tClass.GetCustomAttribute(typeof(REST));
                if (t != null) {
                    Logger.Log("RestProcessor", $"Found REST class {tClass.Name}");
                    REST trest = (REST)t;
                    proxies.Add(tClass.Name, new RestProxy(tClass, injectables));

                    MethodInfo[] methods = tClass.GetMethods();
                    foreach (var methodInfo in methods) {
                        foreach (Type rt in RestTypes) {
                            Attribute rta = methodInfo.GetCustomAttribute(rt);
                            if (rta != null) {
                                RestCall restCall = new RestCall();
                                try {
                                    restCall.className = tClass.Name;
                                    restCall.methodName = methodInfo.Name;
                                    restCall.method = (IHttpMethod)rta;
                                    restCall.baseRest = trest;

                                    Logger.Log("RestProcessor", $"Registering method {methodInfo.Name} for {restCall.method.Method} {trest.Path}{restCall.method.Path}");

                                    AddEndpoint(restCall);
                                } catch (DuplicateRestMethodException) {
                                    Logger.Log("RestProcessor", $"DuplicateRestMethodException: There is already a {restCall.method.Method} {trest.Path}{restCall.method.Path} registered.");
                                }
                            }
                        }
                    }
                }

                // Search for RestExceptionHandler Attribute
                t = tClass.GetCustomAttribute(typeof(RestExceptionHandler));
                if (t != null) {
                    Logger.Log("RestProcessor", $"Found a RestExceptionHandler {tClass.Name}");
                    if (typeof(IRestExceptionHandler).IsAssignableFrom(tClass)) {
                        RestExceptionHandler reh = (RestExceptionHandler)t;
                        if (typeof(Exception).IsAssignableFrom(reh.exceptionType)) {
                            IRestExceptionHandler handler = (IRestExceptionHandler)Activator.CreateInstance(tClass);
                            exceptionHandlers.Add(reh.ExceptionType.Name, handler);
                            Logger.Log("RestProcessor", $"     Registered a custom exception handler for exception \"{reh.ExceptionType.Name}\" for class {tClass.Name}");
                        } else {
                            Logger.Log("RestProcessor", $"     Class {tClass.Name} contains the \"RestExceptionHandler\" attribute the passed type does not inherit Exception class. Skipping it.");
                        }
                    } else {
                        Logger.Log("RestProcessor", $"     Class {tClass.Name} contains the \"RestExceptionHandler\" attribute but does not implement IRestExceptionHandler. Skipping it.");
                    }
                }
            }
            if (proxies.Count != 0 || endpoints.Keys.Count != 0 || exceptionHandlers.Keys.Count != 0) {
                Logger.Log("RestProcessor", $"Initialized {proxies.Count} REST proxies.");
                Logger.Log("RestProcessor", $"Initialized {endpoints.Keys.Count} REST endpoints.");
                Logger.Log("RestProcessor", $"Initialized {exceptionHandlers.Keys.Count} Custom Exception Handlers");
            }
        }

        RestCall GetEndPoint(string path, string method) {
            if (endpoints.ContainsKey(path) && endpoints[path].ContainsKey(method)) {
                return endpoints[path][method];
            }

            return null;
        }

        public IRestExceptionHandler GetExceptionHandler(string exceptionName) {
            if (exceptionHandlers.ContainsKey(exceptionName)) {
                return exceptionHandlers[exceptionName];
            }
            return null;
        }

        public RestResult CallEndPoint(string path, string method, RestRequest request) {
            RestCall rc = GetEndPoint(path, method);
            if (proxies.ContainsKey(rc.className)) {
                object ret = proxies[rc.className].CallMethod(rc.methodName, request);

                if (ret is string) {
                    return new RestResult((string)ret);
                }
                return new RestResult(JsonConvert.SerializeObject(ret), MimeTypes.JSON);

            }
            return new RestResult(new ErrorObject {
                ErrorCode = ErrorCodes.NotFound,
                Message = "Endpoint not found",
                ErrorField = "url"
            }.ToJSON(), MimeTypes.JSON, HttpStatusCode.NotFound);
        }

        public bool ContainsEndPoint(string path, string method) {
            return endpoints.ContainsKey(path) && endpoints[path].ContainsKey(method);
        }

        void AddEndpoint(RestCall restCall) {
            string path = restCall.baseRest.Path + restCall.method.Path;
            if (!endpoints.ContainsKey(path)) {
                endpoints.Add(path, new Dictionary<string, RestCall>());
            }

            if (endpoints[path].ContainsKey(restCall.method.Method)) {
                throw new DuplicateRestMethodException();
            }

            endpoints[path][restCall.method.Method] = restCall;
        }
    }
}
