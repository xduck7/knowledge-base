const WS_URL = "ws://localhost:8080/ws";

class SimpleWsClient {
    private socket: WebSocket | null = null;

    connect() {
        this.socket = new WebSocket(WS_URL);

        this.socket.onopen = () => {
            console.log("WS connected");
            // Отправляем тестовое сообщение
            this.send("hello from TS client");
        };

        this.socket.onmessage = (event: MessageEvent) => {
            console.log("message from server:", event.data);
        };

        this.socket.onerror = (event: Event) => {
            console.error("WS error:", event);
        };

        this.socket.onclose = (event: CloseEvent) => {
            console.log("WS closed:", event.code, event.reason);
        };
    }

    send(message: string) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(message);
        } else {
            console.warn("WS not open, cannot send");
        }
    }

    disconnect() {
        if (this.socket) {
            this.socket.close(1000, "client disconnect");
            this.socket = null;
        }
    }
}

const client = new SimpleWsClient();
client.connect();