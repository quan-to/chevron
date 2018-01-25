/**
 *  Made by Lucas Teske - Original Project: https://github.com/opensatelliteproject/
 */
using System;

namespace RemoteSigner {
    public enum ConsoleMessagePriority {
        INFO,
        WARN,
        ERROR,
        DEBUG
    }

    public class ConsoleMessage : ICloneable {

        public DateTime TimeStamp { get; set; }
        public string Message { get; set; }
        public ConsoleMessagePriority Priority { get; set; }
        public ConsoleMessage(ConsoleMessagePriority priority, string message) {
            TimeStamp = DateTime.Now;
            Message = message;
            Priority = priority;
        }

        public override string ToString() {
            return String.Format("{0}/{1,-5} {2}", TimeStamp.ToLongTimeString(), Priority.ToString(), Message);
        }

        #region ICloneable implementation

        public object Clone() {
            ConsoleMessage cm = new ConsoleMessage(Priority, Message) {
                TimeStamp = TimeStamp
            };
            return cm;
        }

        #endregion
    }
}