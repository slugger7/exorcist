/** @import { LoginResult } from "../types" */
import { server } from "../env.js";
import { fetch } from "./fetch.js";

/**
 * @param {string} username
 * @param {string} password
 * @returns {Promise<LoginResult>}
 */
export const login = async (username, password) => {
  const res = await fetch(`${server()}/login`, {
    method: "POST",
    body: JSON.stringify({ username, password }),
  });

  if (res.status != 201) {
    throw new Error(`could not log in ${username} with given password`);
  }

  if (res.status == 201) {
    const data = await res.json();

    localStorage.setItem("userId", data.userId);
    localStorage.setItem("username", data.username)
    return data;
  }
};

export const logout = async () => {
  const res = await fetch(`${server()}/logout`);

  if (res.status != 200) {
    throw new Error("could not log out");
  }
};
