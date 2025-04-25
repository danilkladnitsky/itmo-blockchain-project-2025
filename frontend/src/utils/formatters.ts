export const formatHash = (hash: string): string => {
    return `${hash.slice(0, 6)}...${hash.slice(-4)}`;
};

export const formatAddress = (address: string): string => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
};

export const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleString();
};

export const formatEther = (value: string): string => {
    const etherValue = Number(value) / 1e18;
    return etherValue.toFixed(4);
};

export const formatGasPrice = (gasPrice: string): string => {
    const gweiValue = Number(gasPrice) / 1e9;
    return gweiValue.toFixed(2);
};
