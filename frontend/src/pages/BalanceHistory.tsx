import {
    AccountBalance as AccountBalanceIcon,
    Timeline as TimelineIcon,
    TrendingDown as TrendingDownIcon,
    TrendingUp as TrendingUpIcon,
} from '@mui/icons-material';
import {
    Box,
    Card,
    CardContent,
    Chip,
    CircularProgress,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
} from '@mui/material';
import React, { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../services/api';

interface BalanceHistory {
  id: number;
  user_id: number;
  amount: number;
  previous_amount: number;
  change_amount: number;
  change_type: 'credit' | 'debit' | 'transfer_in' | 'transfer_out';
  created_at: string;
}

const BalanceHistory: React.FC = () => {
  const { user } = useAuth();
  const [history, setHistory] = useState<BalanceHistory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchBalanceHistory();
  }, [user]);

  const fetchBalanceHistory = React.useCallback(async () => {
    if (!user) return;

    try {
      setLoading(true);
      const response = await api.get(`/api/v1/balances/historical?user_id=${user.id}`);
      setHistory(response.data || []);
    } catch (error: any) {
      setError(error.response?.data || 'Bakiye geçmişi yüklenirken bir hata oluştu');
    } finally {
      setLoading(false);
    }
  }, [user]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('tr-TR', {
      style: 'currency',
      currency: 'TRY',
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('tr-TR');
  };

  const getChangeTypeText = (type: string) => {
    switch (type) {
      case 'credit':
        return 'Para Yatırma';
      case 'debit':
        return 'Para Çekme';
      case 'transfer_in':
        return 'Transfer Gelen';
      case 'transfer_out':
        return 'Transfer Giden';
      default:
        return type;
    }
  };

  const getChangeTypeColor = (type: string) => {
    switch (type) {
      case 'credit':
      case 'transfer_in':
        return 'success';
      case 'debit':
      case 'transfer_out':
        return 'error';
      default:
        return 'default';
    }
  };

  const getChangeIcon = (type: string) => {
    switch (type) {
      case 'credit':
      case 'transfer_in':
        return <TrendingUpIcon color="success" />;
      case 'debit':
      case 'transfer_out':
        return <TrendingDownIcon color="error" />;
      default:
        return <TimelineIcon />;
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" p={3}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        <AccountBalanceIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
        Bakiye Geçmişi
      </Typography>

      {error && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Typography color="error">{error}</Typography>
          </CardContent>
        </Card>
      )}

      <Box>
        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              Bakiye Değişim Geçmişi
            </Typography>
            
            {history.length === 0 ? (
              <Box textAlign="center" p={3}>
                <Typography variant="body2" color="text.secondary">
                  Henüz bakiye değişimi bulunmuyor
                </Typography>
              </Box>
            ) : (
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Tarih</TableCell>
                      <TableCell>İşlem Türü</TableCell>
                      <TableCell align="right">Önceki Bakiye</TableCell>
                      <TableCell align="right">Değişim</TableCell>
                      <TableCell align="right">Yeni Bakiye</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {history.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell>
                          {formatDate(item.created_at)}
                        </TableCell>
                        <TableCell>
                          <Box display="flex" alignItems="center" gap={1}>
                            {getChangeIcon(item.change_type)}
                            <Chip
                              label={getChangeTypeText(item.change_type)}
                              color={getChangeTypeColor(item.change_type) as any}
                              size="small"
                            />
                          </Box>
                        </TableCell>
                        <TableCell align="right">
                          {formatCurrency(item.previous_amount)}
                        </TableCell>
                        <TableCell align="right">
                          <Typography
                            color={item.change_amount >= 0 ? 'success.main' : 'error.main'}
                            fontWeight="bold"
                          >
                            {item.change_amount >= 0 ? '+' : ''}{formatCurrency(item.change_amount)}
                          </Typography>
                        </TableCell>
                        <TableCell align="right">
                          <Typography fontWeight="bold">
                            {formatCurrency(item.amount)}
                          </Typography>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CardContent>
        </Card>
      </Box>
    </Box>
  );
};

export default BalanceHistory; 