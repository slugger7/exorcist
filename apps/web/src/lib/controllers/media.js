/** @import { PageDTO, MediaDTO, MediaOverviewDTO, ProgressDTO, MediaUpdateDTO, MediaUpdatedDTO } from "../../dto" */
/** @import {Ordinal} from "../types" */
import { server } from '../env.js'
import { fetch } from './fetch.js'

/**
 * @param {number} [page]
 * @param {number} [limit]
 * @param {string} [search]
 * @param {string} [orderBy]
 * @param {boolean} [ascending]
 * @returns {Promise<PageDTO<MediaOverviewDTO>>}
 */
export const getAll = async (page = 1, limit = 48, search = "", orderBy = "", ascending = true) => {
  const params = new URLSearchParams()
  params.set("skip", ((page - 1) * limit).toString())
  params.set("limit", limit.toString())
  params.set("asc", ascending.toString())

  if (search !== "") {
    params.set("search", search)
  }

  if (orderBy !== "") {
    params.set("orderBy", orderBy)
  }

  return getAllWithParams(params)
}

/**
 * @param {URLSearchParams} [params]
 * @returns {Promise<PageDTO<MediaOverviewDTO>>}
 */
export const getAllWithParams = async (params = new URLSearchParams()) => {
  const res = await fetch(`${server()}/media?${params.toString()}`)

  return await res.json()
}

/**
 * @param {string} id 
 * @returns {string}
 */
export const videoUrlById = (id) => `${server()}/videos/${id}`

/**
 * @param {string} id
 * @returns {Promise<MediaDTO>}
 */
export const get = async (id) => {
  const res = await fetch(
    `${server()}/media/${id}`
  )

  return await res.json()
}

/**
 * @param {string} id
 * @param {boolean} [physical]
 * @returns {Promise}
 */
export const deleteMedia = async (id, physical = false) => {
  const params = new URLSearchParams()
  params.set("physical", physical.toString())

  const res = await fetch(`${server()}/media/${id}?${params.toString()}`, { method: "DELETE" })

  if (!res.ok) {
    throw new Error("non ok response from server when deleting media")
  }
}

/**
 * @param {string} id
 * @param {number} timestamp
 * @param {boolean} [overwrite]
 * @returns {Promise<ProgressDTO>}
 */
export const updateProgress = async (id, timestamp, overwrite = false) => {
  const params = new URLSearchParams()

  params.set("progress", timestamp.toString())
  params.set("overwrite", overwrite.toString())

  const res = await fetch(`${server()}/videos/${id}?${params.toString()}`, { method: "PUT" })

  return await res.json()
}

/**
 * @param {string} id
 * @param {MediaUpdateDTO} updateDto
 * @returns {Promise<MediaUpdatedDTO>}
 */
export const updateMedia = async (id, updateDto) => {
  const res = await fetch(`${server()}/media/${id}`, { method: "PUT", body: JSON.stringify(updateDto) })

  return await res.json()
}


/** @type {Ordinal[]} */
export const ordinals = [
  {
    text: "Added",
    value: "added"
  },
  {
    text: "Created",
    value: "created"
  },
  {
    text: "Modified",
    value: "modified"
  },
  {
    text: "Title",
    value: "title"
  },
  {
    text: "Size",
    value: "size"
  },
  {
    text: "Runtime",
    value: "runtime"
  },
]