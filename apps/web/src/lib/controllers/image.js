import { server } from "../env";

/**
 * @param {string} id 
 * @returns {string}
 */
export const imageUrlById = (id) => `${server()}/images/${id}`