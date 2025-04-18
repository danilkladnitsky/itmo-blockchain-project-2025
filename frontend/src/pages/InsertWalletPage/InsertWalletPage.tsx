import { Text, TextInput, Button, Icon } from '@gravity-ui/uikit'
import bem from 'bem-cn-lite';

import './styles.scss';
import { CreditCard } from '@gravity-ui/icons';
import { useAppContext } from '@/App.context';

const b = bem('insert-wallet-page');

export const InsertWalletPage = () => {
  const { walletAddress, setWalletAddress, searchWallet } = useAppContext();

  const handleSearch = () => {
    searchWallet();
  }

  return (
    <div className={b()}>
      <div className={b('form')}>
        <Text variant='body-3' className={b('form-title')}>Provide existing wallet address <Icon size={20} data={CreditCard} /></Text>
        <div className={b('form-input')}>
          <TextInput placeholder='Wallet address in hex format' size='xl' hasClear value={walletAddress} onChange={(e) => setWalletAddress(e.target.value)} />
          <Button size='xl' view='action' onClick={handleSearch}>Search</Button>
        </div>
      </div>
    </div>
  )
}
