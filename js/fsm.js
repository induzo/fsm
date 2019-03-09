export const getOutcomeActions = (graph, status, role) => {
  const fromSrc = graph
    .filter(
      s =>
        s.sources && s.sources.filter(src => src.status === status).length > 0
    )
    .map(s => s.sources)
    .reduce((acc, curr) => {
      acc = acc.concat(curr);
      return acc;
    }, [])
    .filter(srx => srx.status === status && srx.roles.includes(role) > 0)
    .map(src => src.action)
    .reduce((acc, curr) => {
      acc.push(curr);
      return acc;
    }, []);

  return fromSrc;
};
