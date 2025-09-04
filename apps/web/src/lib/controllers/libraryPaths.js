/** @import {LibraryPathDTO} from "../../dto" */
import { server } from "../env"
import { fetch } from "./fetch"

/**
 * @param {string} libraryId
 * @param {string} path
 * @returns {Promise<LibraryPathDTO>}
 */
export const create = async (libraryId, path) => {
  const res = await fetch(
    `${server()}/libraryPaths`,
    { method: "POST", body: JSON.stringify({ libraryId, path }) }
  )

  return await res.json()
}

/**
 * @param {string} id 
 * @returns {Promise<LibraryPathDTO>}
 */
export const get = async (id) => {
  const res = await fetch(
    `${server()}/libraryPaths/${id}`
  )

  return await res.json()
}
