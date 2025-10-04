import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { Wallet, Users, BarChart3, LogOut, Menu, X } from 'lucide-react'
import { useState } from 'react'
import NotificationCenter from './NotificationCenter'

const Navbar = () => {
  const { user, logout } = useAuth()
  const location = useLocation()
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: BarChart3 },
    { name: 'Groups', href: '/groups', icon: Users },
   