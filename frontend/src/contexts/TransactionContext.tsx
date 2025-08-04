import React, { createContext, ReactNode, useContext, useEffect, useState } from 'react';
import { api } from '../services/api';
import { useAuth } from './AuthContext';

interface Transaction {
  id: number;
  from_user_id: number;
  to_user_id?: number;
  amount: number;
  type: 'credit' | 'debit' | 'transfer';
  status: 'pending' | 'completed' | 'failed';
  created_at: string;
}

interface Balance {
  user_id: number;
  amount: number;
  last_updated_at: string;
}

interface TransactionContextType {
  transactions: Transaction[];
  balance: Balance | null;
  loading: boolean;
  refreshTransactions: () => Promise<void>;
  refreshBalance: () => Promise<void>;
  performTransaction: (type: 'credit' | 'debit' | 'transfer', data: any) => Promise<void>;
}

const TransactionContext = createContext<TransactionContextType | undefined>(undefined);

export const useTransactions = () => {
  const context = useContext(TransactionContext);
  if (context === undefined) {
    throw new Error('useTransactions must be used within a TransactionProvider');
  }
  return context;
};

interface TransactionProviderProps {
  children: ReactNode;
}

export const TransactionProvider: React.FC<TransactionProviderProps> = ({ children }) => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [balance, setBalance] = useState<Balance | null>(null);
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();

  const refreshTransactions = async () => {
    if (!user) return;
    
    try {
      setLoading(true);
      const response = await api.get(`/api/v1/transactions/history?user_id=${user.id}`);
      setTransactions(response.data || []);
    } catch (error) {
      console.error('Error fetching transactions:', error);
    } finally {
      setLoading(false);
    }
  };

  const refreshBalance = async () => {
    if (!user) return;
    
    try {
      setLoading(true);
      const response = await api.get(`/api/v1/balances/current?user_id=${user.id}`);
      setBalance(response.data);
    } catch (error) {
      console.error('Error fetching balance:', error);
    } finally {
      setLoading(false);
    }
  };

  const performTransaction = async (type: 'credit' | 'debit' | 'transfer', data: any) => {
    if (!user) return;
    
    try {
      setLoading(true);
      
      let endpoint = '';
      let payload = {};
      
      switch (type) {
        case 'credit':
          endpoint = '/api/v1/transactions/credit';
          payload = {
            user_id: user.id,
            amount: data.amount,
          };
          break;
        case 'debit':
          endpoint = '/api/v1/transactions/debit';
          payload = {
            user_id: user.id,
            amount: data.amount,
          };
          break;
        case 'transfer':
          endpoint = '/api/v1/transactions/transfer';
          payload = {
            from_user_id: user.id,
            to_user_id: data.toUserId,
            amount: data.amount,
          };
          break;
      }
      
      await api.post(endpoint, payload);
      
      // Refresh data after successful transaction
      await Promise.all([refreshBalance(), refreshTransactions()]);
    } catch (error) {
      console.error('Transaction error:', error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user) {
      refreshBalance();
      refreshTransactions();
    }
  }, [user]);

  const value: TransactionContextType = {
    transactions,
    balance,
    loading,
    refreshTransactions,
    refreshBalance,
    performTransaction,
  };

  return (
    <TransactionContext.Provider value={value}>
      {children}
    </TransactionContext.Provider>
  );
}; 