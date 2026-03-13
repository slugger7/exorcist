import { server } from "../env";

/**
 * @param {string} id 
 * @returns {string}
 */
export const imageUrlById = (id) => `${server()}/images/${id}`

/**
 * @param {string} mediaId
 * @returns {string}
 */
export const thumbnailUrl = (mediaId) => `${server()}/media/${mediaId}/thumbnail`