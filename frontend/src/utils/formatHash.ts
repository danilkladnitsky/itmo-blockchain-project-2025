export const formatHash = (hash: string) => {
    return hash.slice(0, 7) + '...' + hash.slice(-4);
};
