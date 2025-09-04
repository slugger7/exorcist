import { PING, PONG } from "../constants/websocket"
import { ws, WS_ACTIVE } from "../env"
import { tryParseIntOrDefault } from "../parsing"
import { wsState } from "../state/wsState.svelte"

const heartbeat = tryParseIntOrDefault(import.meta.env.VITE_WEBSOCKET_HEARTBEAT, 30000)
const pingTime = heartbeat * 9 / 10
const pongTime = heartbeat

/** @type {WebSocket} */
let conn
let pingInterval
let pongTimeout

window.addEventListener("beforeunload", () => {
	if (conn) {
		conn.close()
	}
	clearInterval(pingInterval)
	clearTimeout(pongTimeout)
})

export const setupWebsocket = () => {
	const connect = () => {
		if (!WS_ACTIVE()) { 
			console.log("No websockets running")
			return 
		}
		try {
			conn = new WebSocket(`${ws()}/ws`)

			conn.addEventListener("open", onOpen)
			conn.addEventListener("close", onClose)
			conn.addEventListener("message", onMessage)
			conn.addEventListener("error", onError)
		} catch (e) {
			console.error("could not open websocket to server", e)
			conn == null
		}
	}

	const onClose = (evt) => {
		console.warn("Websocket connection closed")
		clearInterval(pingInterval)
		clearTimeout(pongTimeout)
		wsState.active = false
		wsState.connection = null
		connect()
	}

	const onMessage = (evt) => {
		if (evt.data === PONG) {
			console.debug(PONG)
			pong()
			return
		}
	}

	const onOpen = (evt) => {
		console.info("Websoceket opened")
		wsState.active = true;
		wsState.connection = conn
		ping()
		pong()
	}

	const onError = (evt) => {
		console.error(evt)
	}

	const ping = () => {
		pingInterval = setInterval(() => {
			console.debug(PING)
			if (conn && conn.readyState !== WebSocket.CLOSED) {
				try {
					conn.send(PING)
				} catch (e) {
					console.error("error sending ping on websocket", e)
				}
			}
		}, pingTime)
	}

	const pong = () => {
		clearTimeout(pongTimeout)
		pongTimeout = setTimeout(() => {
			console.warn("Timeout reached, recreating websocket if possible")
			clearInterval(pingInterval)
			clearTimeout(pongTimeout)
			conn.close()
			wsState.connection = null
			wsState.active = false
			connect()
		}, pongTime)
	}

	connect()
}