/**
 * @import { MediaOverviewDTO, PageDTO, PersonDTO, PersonOrdinalAllValues } from "../../dto"
 * @import { Ordinal } from "../types"
 */
import { server } from "../env";
import { fetch } from "./fetch";


/**
 * @param {string} [search]
 * @param {PersonOrdinalAllValues} [ordinal]
 * @returns {Promise<PersonDTO[]>} */
export const getAll = async (search = "", ordinal = "count") => {
  const params = new URLSearchParams()

  if (search.length > 0) {
    params.set("search", search)
  }

  params.set("orderBy", ordinal)

  return await getAllWithParams(params)
}

/** 
  * @param {URLSearchParams} [params]
  * @returns {Promise<PersonDTO[]>} */
export const getAllWithParams = async (params = new URLSearchParams()) => {
  const res = await fetch(`${server()}/people?${params.toString()}`)

  return await res.json()
}

export const remove = async (mediaId, personId) => {
  const res = await fetch(`${server()}/media/${mediaId}/people/${personId}`, { method: "DELETE" })

  return res.ok
}

export const add = async (mediaId, personId) => {
  const res = await fetch(`${server()}/media/${mediaId}/people/${personId}`, { method: "PUT" })

  return await res.json()
}

/** @param {string[]} peopleNames */
export const create = async (peopleNames) => {
  const res = await fetch(`${server()}/people`, { method: "POST", body: JSON.stringify(peopleNames) })

  return await res.json()
}

/**
 * @param {string} id 
 * @param {URLSearchParams} [params]
 * @returns {Promise<PageDTO<MediaOverviewDTO>>} */
export const getMediaWithParams = async (id, params = new URLSearchParams()) => {
  const res = await fetch(`${server()}/people/${id}/media?${params.toString()}`)

  return await res.json()
}

/** @type {Ordinal[]} */
export const ordinals = [
  {
    text: "Count",
    value: "count"
  },
  {
    text: "Name",
    value: "name"
  }
]

/**
 * @param {string} id
 * @param {string} name
 * @returns {Promise<PersonDTO>}
 */
export const updatePerson = async (id, name) => {
  const res = await fetch(`${server()}/people/${id}`, { method: "PUT", body: JSON.stringify({ name }) })

  return await res.json()
}