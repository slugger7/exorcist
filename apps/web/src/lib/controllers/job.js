/** 
 * @import { JobDTO, PageDTO } from "../../dto"
 * @import { JobStatusEnum, JobTypeEnum } from "../../dto/model"
 * @import { JobData } from "../types"
 */
import { server } from "../env";
import { fetch } from "./fetch";

/**
 * @param {JobTypeEnum} type 
 * @param {JobData} data 
 * @returns {Promise<JobDTO>}
 */
export const create = async (type, data) => {
  const res = await fetch(`${server()}/jobs`, {
    method: "POST",
    body: JSON.stringify({
      type,
      data
    })
  })

  return await res.json()
}

/**
 * 
 * @param {string} parent 
 * @param {JobStatusEnum[]} statuses 
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
  params.set("limit", limit)
  params.set("skip", (limit * (page - 1)).toString())

  const res = await fetch(`${server()}/jobs?${params.toString()}`)

  return await res.json()
}