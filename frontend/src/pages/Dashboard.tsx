import {
    AccountBalance as AccountBalanceIcon,
    Add as AddIcon,
    AdminPanelSettings as AdminIcon,
    History as HistoryIcon,
    Person as PersonIcon,
    Refresh as RefreshIcon,
    Remove as RemoveIcon,
    SwapHoriz as SwapHorizIcon,
    SwapHoriz as TransferIcon,
    TrendingDown as TrendingDownIcon,
    TrendingUp as TrendingUpIcon,
} from '@mui/icons-material';
import {
    Alert,
    Box,
    Button,
    Card,
    CardContent,
    Chip,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    IconButton,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    TextField,
    Tooltip,
    Typography,
} from '@mui/material';
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useTransactions } from '../contexts/TransactionContext';

interface Transaction {
  id: number;
  from_user_id: number;
  to_user_id?: number;
  amount: number;
  type: 'credit' | 'debit' | 'transfer';
  status: 'pending' | 'completed' | 'failed';
  created_at: string;
}

const Dashboard: React.FC = () => {
  const { balance, transactions, loading, performTransaction, refreshBalance, refreshTransactions } = useTransactions();
  const { user } = useAuth();
  const navigate = useNavigate();
  
  const [dialogOpen, setDialogOpen] = useState(false);
  const [transactionType, setTransactionType] = useState<'credit' | 'debit' | 'transfer'>('credit');
  const [transactionData, setTransactionData] = useState({
    amount: '',
    toUserId: '',
  });
  const [transactionLoading, setTransactionLoading] = useState(false);
  const [error, setError] = useState('');

  const handleOpenDialog = (type: 'credit' | 'debit' | 'transfer') => {
    setTransactionType(type);
    setTransactionData({ amount: '', toUserId: '' });
    setError('');
    setDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setDialogOpen(false);
    setTransactionData({ amount: '', toUserId: '' });
    setError('');
  };

  const handleTransaction = async () => {
    if (!transactionData.amount) {
      setError('Lütfen miktar girin');
      return;
    }

    if (transactionType === 'transfer' && !transactionData.toUserId) {
      setError('Lütfen alıcı kullanıcı ID girin');
      return;
    }

    setTransactionLoading(true);
    setError('');

    try {
      const data = {
        amount: parseFloat(transactionData.amount),
        ...(transactionType === 'transfer' && { toUserId: parseInt(transactionData.toUserId) }),
      };

      await performTransaction(transactionType, data);
      handleCloseDialog();
    } catch (error: any) {
      setError(error.response?.data || 'İşlem sırasında bir hata oluştu');
    } finally {
      setTransactionLoading(false);
    }
  };

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case 'credit':
        return <TrendingUpIcon color="success" />;
      case 'debit':
        return <TrendingDownIcon color="error" />;
      case 'transfer':
        return <SwapHorizIcon color="info" />;
      default:
        return <HistoryIcon />;
    }
  };

  const getTransactionTypeText = (type: string) => {
    switch (type) {
      case 'credit':
        return 'Para Yatırma';
      case 'debit':
        return 'Para Çekme';
      case 'transfer':
        return 'Transfer';
      default:
        return type;
    }
  };

  const getTransactionAmountColor = (type: string) => {
    switch (type) {
      case 'credit':
        return 'success';
      case 'debit':
        return 'error';
      case 'transfer':
        return 'info';
      default:
        return 'default';
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('tr-TR', {
      style: 'currency',
      currency: 'TRY',
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('tr-TR');
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 2, mb: 2 }}>
        {/* Balance Card */}
        <Box sx={{ flex: { xs: '1', md: '0 0 250px' } }}>
          <Card sx={{ height: '100%' }}>
            <CardContent sx={{ p: 2 }}>
              <Box display="flex" alignItems="center" justifyContent="space-between" mb={1.5}>
                <Typography variant="h6" component="h2">
                  <AccountBalanceIcon sx={{ mr: 1, verticalAlign: 'middle', fontSize: '1.2rem' }} />
                  Bakiye
                </Typography>
                <Tooltip title="Yenile">
                  <IconButton onClick={refreshBalance} disabled={loading} size="small">
                    <RefreshIcon />
                  </IconButton>
                </Tooltip>
              </Box>
              <Typography variant="h4" component="div" color="primary" gutterBottom>
                {balance ? formatCurrency(balance.amount) : formatCurrency(0)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Son güncelleme: {balance ? formatDate(balance.last_updated_at) : 'Bilgi yok'}
              </Typography>
            </CardContent>
          </Card>
        </Box>

        {/* Quick Actions */}
        <Box sx={{ flex: '1' }}>
          <Card sx={{ height: '100%' }}>
            <CardContent sx={{ p: 2 }}>
              <Typography variant="h6" component="h2" gutterBottom>
                Hızlı İşlemler
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, gap: 1.5 }}>
                <Button
                  fullWidth
                  variant="contained"
                  color="success"
                  size="medium"
                  startIcon={<AddIcon />}
                  onClick={() => handleOpenDialog('credit')}
                  sx={{ height: 45 }}
                >
                  Para Yatır
                </Button>
                <Button
                  fullWidth
                  variant="contained"
                  color="error"
                  size="medium"
                  startIcon={<RemoveIcon />}
                  onClick={() => handleOpenDialog('debit')}
                  sx={{ height: 45 }}
                >
                  Para Çek
                </Button>
                <Button
                  fullWidth
                  variant="contained"
                  color="info"
                  size="medium"
                  startIcon={<TransferIcon />}
                  onClick={() => handleOpenDialog('transfer')}
                  sx={{ height: 45 }}
                >
                  Transfer
                </Button>
              </Box>
            </CardContent>
          </Card>
        </Box>
      </Box>

      {/* Quick Links */}
      <Box sx={{ mb: 2 }}>
        <Card>
          <CardContent sx={{ p: 2 }}>
            <Typography variant="h6" component="h2" gutterBottom>
              Hızlı Linkler
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, gap: 1.5 }}>
              <Button
                variant="outlined"
                startIcon={<HistoryIcon />}
                onClick={() => navigate('/balance-history')}
                fullWidth
                size="small"
              >
                Bakiye Geçmişi
              </Button>
              <Button
                variant="outlined"
                startIcon={<PersonIcon />}
                onClick={() => navigate('/profile')}
                fullWidth
                size="small"
              >
                Profil Ayarları
              </Button>
              {user?.role === 'admin' && (
                <Button
                  variant="outlined"
                  startIcon={<AdminIcon />}
                  onClick={() => navigate('/admin')}
                  fullWidth
                  size="small"
                >
                  Admin Paneli
                </Button>
              )}
            </Box>
          </CardContent>
        </Card>
      </Box>

      {/* Transaction History */}
      <Card>
        <CardContent sx={{ p: 2 }}>
          <Box display="flex" alignItems="center" justifyContent="space-between" mb={1.5}>
            <Typography variant="h6" component="h2">
              İşlem Geçmişi
            </Typography>
            <Tooltip title="Yenile">
              <IconButton onClick={refreshTransactions} disabled={loading} size="small">
                <RefreshIcon />
              </IconButton>
            </Tooltip>
          </Box>
          
          {loading ? (
            <Box display="flex" justifyContent="center" p={3}>
              <CircularProgress />
            </Box>
          ) : transactions.length === 0 ? (
            <Box textAlign="center" p={3}>
              <Typography variant="body2" color="text.secondary">
                Henüz işlem bulunmuyor
              </Typography>
            </Box>
          ) : (
            <List>
              {transactions.map((transaction: Transaction) => (
                <ListItem key={transaction.id} divider>
                  <ListItemIcon>
                    {getTransactionIcon(transaction.type)}
                  </ListItemIcon>
                  <ListItemText
                    primary={
                      <Box display="flex" justifyContent="space-between" alignItems="center">
                        <Typography variant="body1">
                          {getTransactionTypeText(transaction.type)}
                        </Typography>
                        <Chip
                          label={formatCurrency(transaction.amount)}
                          color={getTransactionAmountColor(transaction.type) as any}
                          size="small"
                        />
                      </Box>
                    }
                    secondary={
                      <Box display="flex" justifyContent="space-between" alignItems="center">
                        <Typography variant="body2" color="text.secondary">
                          {formatDate(transaction.created_at)}
                        </Typography>
                        <Chip
                          label={transaction.status}
                          color={transaction.status === 'completed' ? 'success' : 'warning'}
                          size="small"
                        />
                      </Box>
                    }
                  />
                </ListItem>
              ))}
            </List>
          )}
        </CardContent>
      </Card>

      {/* Transaction Dialog */}
      <Dialog open={dialogOpen} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>
          {transactionType === 'credit' && 'Para Yatır'}
          {transactionType === 'debit' && 'Para Çek'}
          {transactionType === 'transfer' && 'Transfer Et'}
        </DialogTitle>
        <DialogContent sx={{ pt: 1 }}>
          {error && (
            <Alert severity="error" sx={{ mb: 1.5 }}>
              {error}
            </Alert>
          )}
          
          <TextField
            fullWidth
            label="Miktar (₺)"
            type="number"
            value={transactionData.amount}
            onChange={(e) => setTransactionData({ ...transactionData, amount: e.target.value })}
            margin="dense"
            inputProps={{ min: 0.01, step: 0.01 }}
            required
          />
          
          {transactionType === 'transfer' && (
            <TextField
              fullWidth
              label="Alıcı Kullanıcı ID"
              type="number"
              value={transactionData.toUserId}
              onChange={(e) => setTransactionData({ ...transactionData, toUserId: e.target.value })}
              margin="dense"
              inputProps={{ min: 1 }}
              required
            />
          )}
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={handleCloseDialog} disabled={transactionLoading} size="small">
            İptal
          </Button>
          <Button
            onClick={handleTransaction}
            variant="contained"
            disabled={transactionLoading}
            size="small"
            startIcon={transactionLoading ? <CircularProgress size={16} /> : undefined}
          >
            {transactionLoading ? 'İşlem Yapılıyor...' : 'Onayla'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Dashboard; 