import React, { useState, useEffect } from 'react'
import { Settings, UserPlus, CheckCircle, Copy } from 'lucide-react'
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query'
import { groupApi } from '../services/api'
import toast from 'react-hot-toast'
import type { Group, User } from '../types'
import RoundContributions from './RoundContributions'

interface GroupManagementProps {
  group: Group
  currentUser: User
}

const GroupManagement: React.FC<GroupManagementProps> = ({
  group,
  currentUser,
}) => {
  const [showInviteModal, setShowInviteModal] = useState(false)
  const [showApproveModal, setShowApproveModal] = useState(false)
  const [showActivateModal, setShowActivateModal] = useState(false)
  const [inviteEmail, setInviteEmail] = useState('')
  const [availableUsers, setAvailableUsers] = useState<User[]>([])
  const [groupSettings, setGroupSettings] = useState({
    contribution_amount: 0,
    contribution_period: 30,
    payout_order: [] as string[]
  })
  const queryClient = useQueryClient()

  const inviteUserMutation = useMutation({
    mutationFn: (email: string) => groupApi.inviteToGroup(group.ID, { email }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Invitation sent successfully!')
      setShowInviteModal(false)
      setInviteEmail('')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to send invitation')
    }
  })

  const approveGroupMutation = useMutation({
    mutationFn: () => groupApi.approveGroup(group.ID),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Group approved successfully!')
      setShowApproveModal(false)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to approve group')
    }
  })

  const activateGroupMutation = useMutation({
    mutationFn: (settings: any) => groupApi.activateGroup(group.ID, settings),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Group activated successfully!')
      setShowActivateModal(false)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to activate group')
    }
  })

  // Fetch available users when invite modal opens
  useEffect(() => {
    if (showInviteModal) {
      groupApi.getNonGroupMembers(group.ID)
        .then(response => setAvailableUsers(response.data))
        .catch(error => console.error('Failed to fetch available users:', error))
    }
  }, [showInviteModal, group.ID])

  const isAdmin = group.Members?.find(m =>
    m.UserID === currentUser.id && ['creator', 'admin'].includes(m.Role)
  )

  const isCreator = group.Members?.find(m =>
    m.UserID === currentUser.id && m.Role === 'creator'
  )

  const approvedMembers = group.Members?.filter(m => m.Status === 'approved') || []
  const pendingMembers = group.Members?.filter(m => m.Status === 'pending') || []
 