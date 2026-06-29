export const formatNumber = (val: number): string => {
  val = Math.floor(val)
  const accum = []
  while (val > 1000) {
    accum.push(String(val % 1000).padStart(3, '0'))
    val = Math.floor(val / 1000)
  }
  accum.push(val)
  accum.reverse()
  return accum.join(',')
}
