/** @import { MediaOverviewDTO, PageDTO, PlaylistDTO } from "../types" */
import { server } from "../env";
import { fetch } from "./fetch";

/** @returns {Promise<PlaylistDTO[]>}*/
export const get = async () => {
  const res = await fetch(`${server()}/playlists`)

  return await res.json()
}

/**
 * @param {string} id
 * @param {URLSearchParams} [params]
 * @returns {Promise<PageDTO, MediaOverviewDTO>}
 */
export const getMediaWithParams = async (id, params = new URLSearchParams) => {
  const res = await fetch(`${server()}/playlists/${id}/media?${params.toString()}`)

  return await res.json()
}

/** 
 * @param {string} name 
 * @returns {Promise}
*/
export const create = async (name) => {

  const res = await fetch(`${server()}/playlists`, { method: "POST", body: JSON.stringify([{ name }]) })

  return await res.json()
}

/**
 * @param {string[]} playlists
 * @param {string[]} media
 * @returns {Promise}
 */
export const addMedia = async (playlists, media) => {
  await Promise.all(playlists.map(async p => {
    await fetch(`${server()}/playlists/${p}/media`,
      { method: "PUT", body: JSON.stringify(media.map(m => ({ mediaId: m }))) })
  }))
}

/**
 * @param {string} id
 * @param {string} name
 * @returns {Promise<PlaylistDTO>}
 */
export const updatePlaylist = async (id, name) => {
  const res = await fetch(`${server()}/playlists/${id}`, { method: "PUT", body: JSON.stringify({ name }) })

  return await res.json()
}