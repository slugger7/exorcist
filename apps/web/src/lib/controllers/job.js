/** 
 * @import { CreateJobDTO, JobDTO, PageDTO } from "../../dto"
 * @import { JobStatusEnum } from "../../dto/model"
 */
import { server } from "../env";
import { fetch } from "./fetch";

/**
 * @param {CreateJobDTO} data 
 * @returns {Promise<JobDTO>}
 */
export const create = async (data) => {
  const res = await fetch(`${server()}/jobs`, {
    method: "POST",
    body: JSON.stringify(data)
  })

  return await res.json()
}

/**
 * @param {number} page
 * @param {number} limit
 * @param {string} parent 
 * @param {JobStatusEnum[]} statuses 
 * @param {string[]} types
 * @returns {Promise<PageDTO<JobDTO>>}
 */
export const getAll = async (page, limit, parent, statuses = [], types = []) => {
  const params = new URLSearchParams()
  if (parent) {
    params.set("parent", parent)
  }

  statuses.forEach(status => {
    params.set("status", status)
  });
  types.forEach(type => {
    params.set("type", type)
  })
  params.set("limit", limit.toString())
  params.set("skip", (limit * (page - 1)).toString())

  const res = await fetch(`${server()}/jobs?${params.toString()}`)

  return await res.json()
}
