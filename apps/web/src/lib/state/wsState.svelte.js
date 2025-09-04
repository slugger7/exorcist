/** @type {{active: boolean, connection: WebSocket}} */
export const wsState = $state({
  active: false,
  connection: null
})