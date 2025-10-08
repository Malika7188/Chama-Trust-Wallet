import React, { useState } from 'react'
import { Bell, Check, X, Users, DollarSign, UserPlus } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api, { notificationApi } from '../services/api'
import type { Notification, GroupInvitation } from '../types'
import toast from 'react-hot-toast'

interface NotificationCenterProps {
  isCollapsed: boolean;
}

const NotificationCenter: React.FC<NotificationCenterProps> = ({ isCollapsed }) => {
  const [showNotifications, setShowNotifications] = useState(false)
  const [selected, setSelected] = useState<string[]>([])
  const queryClient = useQueryClient()
  const clearNotificationMutation = useMutation({
    mutationFn: (id: string) => notificationApi.clearNotification(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
      toast.success('Notification cleared')
    },
    onError: () => {
      toast.error('Failed to clear notification')
    }
  })

  const clearSelected = () => {
    selected.forEach(id => clearNotificationMutation.mutate(id))
    setSelected([])
  }

  const { data: notifications = [] } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => {
      console.log('🔍 Fetching notifications...')
      return notificationApi.getNotifications().then((res: { data: Notification[] }) => {
        console.log('✅ Notifications received:', res.data)
        return res.data
      })
    },
    refetchInterval: 30000, // Refetch every 30 seconds
  })

  const { data: invitations = [] } = useQuery({
    queryKey: ['invitations'],
    queryFn: () => {
      console.log('🔍 Fetching invitations...')
      return notificationApi.getInvitations().then((res: { data: GroupInvitation[] }) => {
        console.log('✅ Invitations received:', res.data)
        return res.data
      })
    },
    refetchInterval: 30000, // Refetch every 30 seconds
  })

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) => notificationApi.markAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
    }
  })

  const acceptInvitationMutation = useMutation({
    mutationFn: (id: string) => notificationApi.acceptInvitation(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invitations'] })
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      queryClient.invalidateQueries({ queryKey: ['userGroups'] }) // Add this if you have user-specific groups
      toast.success('Invitation accepted! You are now a member of the group.')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to accept invitation')
    }
  })

  const rejectInvitationMutation = useMutation({
    mutationFn: (id: string) => notificationApi.rejectInvitation(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invitations'] })
      toast.success('Invitation rejected')
    }
  })

  const unreadCount = notifications.filter(n => !n.Read).length + invitations.length

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'contribution_reminder':
        return <DollarSign className="w-5 h-5 text-yellow-500" />
      case 'payout_approved':
        return <Check className="w-5 h-5 text-green-500" />
      case 'new_member_request':
        return <UserPlus className="w-5 h-5 text-blue-500" />
      case 'admin_promotion':
        return <Users className="w-5 h-5 text-purple-500" />
      default:
        return <Bell className="w-5 h-5 text-gray-500" />
    }
  }

  const handleMarkAsRead = (id: string) => {
    markAsReadMutation.mutate(id)
  }

  const handleAcceptInvitation = (id: string) => {
    acceptInvitationMutation.mutate(id)
  }

  const handleRejectInvitation = (id: string) => {
    rejectInvitationMutation.mutate(id)
  }

  return (
    <div className="relative">
      <button
        onClick={() => setShowNotifications(!showNotifications)}
        className={`flex items-center w-full px-4 py-3 rounded-lg transition-colors duration-200 text-white hover:bg-[#2ecc71] hover:text-[#1a237e] ${showNotifications ? 'bg-[#2ecc71] text-[#1a237e]' : ''}`}
        style={{ outline: 'none' }}
      >
        <Bell className="h-6 w-6" />
        {!isCollapsed && <span className="font-medium ml-3">Notifications</span>}
        {unreadCount > 0 && (
          <span className={`absolute top-2 right-6 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center ${isCollapsed ? 'right-2' : ''}`}>
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>

      {showNotifications && (
  <div className="absolute left-0 top-12 w-80 backdrop-blur-lg bg-white/80 border border-gray-200 rounded-2xl shadow-2xl z-50 max-h-96 overflow-y-auto transition-all duration-300">
          <div className="p-4 border-b border-gray-200 flex items-center justify-between">
            <h3 className="text-lg font-semibold text-[#1a237e]">Notifications</h3>
            <span className="text-xs text-gray-500">{unreadCount > 0 ? `${unreadCount} unread` : 'All read'}</span>
          </div>

          {/* Invitations */}
          {invitations.length > 0 && (
            <div className="border-b border-gray-200">
            