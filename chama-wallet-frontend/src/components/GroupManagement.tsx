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
  const hasMinimumMembers = approvedMembers.length >= (group.MinMembers || 3)
  const isGroupFull = approvedMembers.length >= (group.MaxMembers || 20)

  const handleInvite = (e: React.FormEvent) => {
    e.preventDefault()
    inviteUserMutation.mutate(inviteEmail)
  }

  const handleApproveGroup = () => {
    approveGroupMutation.mutate()
  }

  const handleActivate = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('Activating group with settings:', groupSettings)

    // Ensure payout order is not empty
    if (groupSettings.payout_order.length === 0) {
      toast.error('Payout order cannot be empty')
      return
    }

    activateGroupMutation.mutate({
      contribution_amount: groupSettings.contribution_amount,
      contribution_period: groupSettings.contribution_period,
      payout_order: groupSettings.payout_order,
    })
  }

  // Initialize payout order with member IDs in default order
  useEffect(() => {
    if (showActivateModal && approvedMembers.length > 0) {
      const defaultPayoutOrder = approvedMembers.map(member => member.UserID)
      console.log('Initializing payout order:', defaultPayoutOrder)
      console.log('Approved members:', approvedMembers.map(m => ({ id: m.UserID, name: m.User.name })))
      
      setGroupSettings(prev => ({
        ...prev,
        contribution_amount: prev.contribution_amount || 0,
        contribution_period: prev.contribution_period || 30,
        payout_order: defaultPayoutOrder
      }))
    }
  }, [showActivateModal, approvedMembers])

  useEffect(() => {
    if (group.PayoutOrder) {
      console.log('Raw PayoutOrder:', group.PayoutOrder)
      try {
        const parsed = JSON.parse(group.PayoutOrder)
        console.log('Parsed PayoutOrder:', parsed)
        console.log('Approved Members:', approvedMembers.map(m => ({ id: m.UserID, name: m.User.name })))
      } catch (error) {
        console.error('Error parsing PayoutOrder:', error)
      }
    }
  }, [group.PayoutOrder, approvedMembers])

  // Add this query to fetch the payout schedule
  const { data: payoutScheduleResponse } = useQuery({
    queryKey: ['payout-schedule', group.ID],
    queryFn: () => groupApi.getPayoutSchedule(group.ID),
    enabled: group.Status === 'active'
  })

  const payoutSchedule = payoutScheduleResponse?.data || []

  return (
    <div className="space-y-6">
      {/* Group Status */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">Group Status</h3>
          <span className={`px-3 py-1 rounded-full text-sm font-medium ${group.Status === 'active' ? 'bg-green-100 text-green-800' :
              group.Status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                'bg-gray-100 text-gray-800'
            }`}>
          