/**
 * @param {number} runtime
 * @returns {string}
 */
export const formatRuntime = (runtime) => {
  const f = 60

  const hours = Math.trunc(runtime / (f * f))
  const minutes = Math.trunc((runtime - (hours * f * f)) / (f))
  const seconds = Math.round(runtime - (hours * f * f) - (minutes * f))

  return `${hours}:${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`
}