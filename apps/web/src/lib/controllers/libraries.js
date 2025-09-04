/** @import { LibraryDTO, LibraryPathDTO, PageDTO, MediaOverviewDTO } from "../../dto" */
import { server } from "../env.js";
import { fetch } from "./fetch.js";

/**
 * @returns {Promise<[LibraryDTO]>}
 */
export const getLibraries = async () => {
  const res = await fetch(`${server()}/libraries`);

  if (res.status != 200) {
    console.log("could not get routes");
  }

  return await res.json()
};

/**
 * @param {string} name
 * @returns {Promise<LibraryDTO>}
 */
export const create = async (name) => {
  const res = await fetch(
    `${server()}/libraries`,
    { method: "POST", body: JSON.stringify({ name }) }
  );

  return await res.json()
}

/**
 * @param {string} id
 * @returns {Promise<[LibraryPathDTO]>}
 */
export const getPaths = async (id) => {
  const res = await fetch(
    `${server()}/libraries/${id}/libraryPaths`
  )

  return await res.json()
}

/**
 * @param {string} id 
 * @param {URLSearchParams} [params]
 * @returns {Promise<PageDTO<MediaOverviewDTO>>}
 */
export const getMediaWithParams = async (id, params = new URLSearchParams) => {
  const res = await fetch(`${server()}/libraries/${id}/media?${params.toString()}`)

  return await res.json()
}

/**
 * @param {string} id
 * @param {string} name
 * @returns {Promise<LibraryDTO>}
 */
export const updateLibrary = async (id, name) => {
  const res = await fetch(`${server()}/libraries/${id}`, { method: "PUT", body: JSON.stringify({ name }) })

  return await res.json()
}