export const getWalletAddress = (address: string, showFull = false) => {
    if (showFull) {
        return address;
    }

    const leftPart = address.slice(0, 7);
    const rightPart = address.slice(-4);
    return `${leftPart}...${rightPart}`;
};
