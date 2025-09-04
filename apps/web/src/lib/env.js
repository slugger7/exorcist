export const server = () => {
  let url = ""
  if (import.meta.env.DEV) {
    url = `${import.meta.env.VITE_SERVER_PROTOCOL}://${import.meta.env.VITE_SERVER_HOSTNAME}:${import.meta.env.VITE_SERVER_PORT}`;
  }

  if (import.meta.env.PROD) {
    url = `${window.location.protocol}//${window.location.host}`
  }

  return `${url}/api`
}

export const ws = () => {
  let url = ""
  if (import.meta.env.DEV) {
    url = `${import.meta.env.VITE_SERVER_WEBSOCKET_PROTOCOL}://${import.meta.env.VITE_SERVER_HOSTNAME}:${import.meta.env.VITE_SERVER_PORT}`
  }

  // need to figure out if it is https or http and change the websocket accordingly
  if (import.meta.env.PROD) {
    url = `ws://${window.location.host}`
  }

  return `${url}/api`
}

export const WS_ACTIVE = () => import.meta.env.VITE_WS.toLowerCase() === "true"