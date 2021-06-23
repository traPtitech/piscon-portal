const sidebarOpened = state => state.app.sidebar.opened
const toggleWithoutAnimation = state => state.app.sidebar.withoutAnimation
const config = state => state.app.config
const palette = state => state.app.config.palette
const isLoading = state => state.app.isLoading
const maxScore = state => (state.Team.results || []).reduce((a, b) => { return a.score < b.score ? b : a }, {score: 0})
const resultCount = state => state.AllResults.reduce((a, b) => a + (b.results || []).length, 0)
const rankingData = state => state.AllResults.map(team => {
  const res = {}
  res.name = team.name
  res.result = (team.results || []).filter(result => result.pass).reduce((a, b) => { return a.score < b.score ? b : a }, {score: 0})
  return res
}).sort((a, b) => b.result.score - a.result.score)
const lastResult = state => {
  const l = state.Team.results.length
  return l > 0 ? JSON.stringify(state.Team.results[l - 1], null, '  ') : 'まだベンチマークは行われていません'
}

const recentResults = state => {
  const results = state.AllResults.reduce((a, b) => a.concat((b.results || [])), []).sort((a, b) => b.id - a.id)
  if (results.length > 20) {
    return results.splice(0, 20)
  }
  return results
}

export {
  toggleWithoutAnimation,
  sidebarOpened,
  config,
  palette,
  isLoading,
  maxScore,
  resultCount,
  rankingData,
  lastResult,
  recentResults
}
