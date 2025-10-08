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
   