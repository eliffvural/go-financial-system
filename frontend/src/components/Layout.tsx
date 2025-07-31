import {
    AccountBalance as AccountBalanceIcon,
    AdminPanelSettings as AdminIcon,
    Dashboard as DashboardIcon,
    History as HistoryIcon,
    Logout as LogoutIcon,
    Person as PersonIcon,
} from '@mui/icons-material';
import {
    AppBar,
    Avatar,
    Box,
    Button,
    Container,
    IconButton,
    Menu,
    MenuItem,
    Toolbar,
    Typography
} from '@mui/material';
import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
    handleClose();
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBar
        position="static"
        sx={{
          background: '#ffffff',
          boxShadow: '0 2px 4px -1px rgba(0, 0, 0, 0.1)',
        }}
      >
        <Toolbar>
          <Box display="flex" alignItems="center" sx={{ flexGrow: 1 }}>
            <AccountBalanceIcon sx={{ mr: 2, color: 'primary.main' }} />
            <Typography variant="h6" component="div" sx={{ color: 'text.primary' }}>
              Go Financial System
            </Typography>
          </Box>

          <Box display="flex" alignItems="center" gap={2}>
            <Typography variant="body2" sx={{ color: 'text.primary', fontWeight: 500 }}>
              Hoş geldin, {user?.username}
            </Typography>
            
            <Button
              color="primary"
              startIcon={<DashboardIcon />}
              onClick={() => navigate('/dashboard')}
              sx={{ textTransform: 'none', color: 'text.primary' }}
            >
              Dashboard
            </Button>
            
            <Button
              color="primary"
              startIcon={<PersonIcon />}
              onClick={() => navigate('/profile')}
              sx={{ textTransform: 'none', color: 'text.primary' }}
            >
              Profil
            </Button>
            
            <Button
              color="primary"
              startIcon={<HistoryIcon />}
              onClick={() => navigate('/balance-history')}
              sx={{ textTransform: 'none', color: 'text.primary' }}
            >
              Bakiye Geçmişi
            </Button>
            
            {user?.role === 'admin' && (
              <Button
                color="primary"
                startIcon={<AdminIcon />}
                onClick={() => navigate('/admin')}
                sx={{ textTransform: 'none', color: 'text.primary' }}
              >
                Admin
              </Button>
            )}
            
            <IconButton
              size="large"
              onClick={handleMenu}
              sx={{ color: 'primary.main' }}
            >
              <Avatar sx={{ width: 32, height: 32, bgcolor: 'primary.main' }}>
                <PersonIcon />
              </Avatar>
            </IconButton>
            
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleClose}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'right',
              }}
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
            >
              <MenuItem onClick={handleLogout}>
                <LogoutIcon sx={{ mr: 1 }} />
                Çıkış Yap
              </MenuItem>
            </Menu>
          </Box>
        </Toolbar>
      </AppBar>

      <Container maxWidth="xl" sx={{ flexGrow: 1, py: 2 }}>
        {children}
      </Container>
    </Box>
  );
};

export default Layout; 