import { createContext, useContext, useState } from "react";

type AppContextType = {
    walletAddress: string;
    setWalletAddress: (address: string) => void;
    searchWallet: () => void;
}

const AppContext = createContext<AppContextType>({} as AppContextType);

export const useAppContext = () => {
    return useContext(AppContext);
}

export const AppContextProvider = ({ children }: { children: React.ReactNode }) => {
    const [walletAddress, setWalletAddress] = useState('');

    const searchWallet = async () => {
        const response = await fetch(`https://api.etherscan.io/api?module=account&action=balance&address=${walletAddress}&tag=latest&apikey=YourApiKeyToken`);
        const data = await response.json();
        console.log(data);
    }

    return (
        <AppContext.Provider value={{ walletAddress, setWalletAddress, searchWallet }}>{children}</AppContext.Provider>
    );
}