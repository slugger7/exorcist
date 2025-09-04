/** @import { ResetPasswordDTO } from "../../dto" */
import { server } from "../env";
import { fetch } from "./fetch";

/**
 * 
 * @param {ResetPasswordDTO} updatePasswordDTO 
 */
export const updatePassword = async (updatePasswordDTO) => {
  const res = await fetch(`${server()}/users`, {
    method: "PUT",
    body: JSON.stringify(updatePasswordDTO)
  })

  return await res.json()
}

/**
 * @param {string} mediaId
 */
export const removeFavourite = async (mediaId) => {
  const res = await fetch(`${server()}/users/favourites/${mediaId}`, { method: "DELETE" })

  return res.ok
}

/**
 * @param {string} mediaId
 */
export const addFavourite = async (mediaId) => {
  const res = await fetch(`${server()}/users/favourites/${mediaId}`, { method: "PUT" })

  return res.ok
}