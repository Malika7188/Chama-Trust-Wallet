import React, { useState } from 'react'
import { DollarSign, Check, X, Clock, Users } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { groupApi, payoutApi } from '../services/api'
import type { Group, User, PayoutRequest } from '../types'
import toast from 'react-hot-toast'

interface PayoutManagementProps {
  group: Group
  currentUser: User
}

const PayoutManagement: React.FC<PayoutManagementProps> = ({ group, currentUser }) => {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [payoutData, setPayoutData] = useState({
    recipient_id: '',
    amount: 0,
    round: 1
  })
  const queryClient = useQueryClient()

  const { data: payoutRequestsResponse } = useQuery({
    queryKey: ['payout-requests', group.ID],
    queryFn: () => groupApi.getPayoutRequests(group.ID)
  })

  const payoutRequests = payoutRequestsResponse?.data || []

  const createPayoutMutation = useMutation({
    mutationFn: (data: any) => groupApi.createPayoutRequest(group.ID, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payout-requests', group.ID] })
      toast.success('Payout request created!')
      setShowCreateModal(false)
      setPayoutData({ recipient_id: '', amount: 0, round: 1 })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to create payout request')
    }
  })

  const approvePayoutMutation = useMutation({
    mutationFn: ({ id, approved }: { id: string, approved: boolean }) => 
      payoutApi.approvePayout(id, { approved }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payout-requests', group.ID] })
      toast.success('Payout decision recorded!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to process payout decision')
    }
  })

  const currentUserMember = group.Members?.find(m => m.UserID === currentUser.id)
  const isAdmin = currentUserMember && ['creator', 'admin'].includes(currentUserMember.Role)
  const approvedMembers = group.Members?.filter(m => m.Status === 'approved') || []

  const handleCreatePayout = (e: React.FormEvent) => {
    e.preventDefault()
    createPayoutMutation.mutate(payoutData)
  }

  const handleApprovePayout = (payoutId: string, approved: boolean) => {
    approvePayoutMutation.mutate({ id: payoutId, approved })
  }

  const getPayoutStatus = (payout: PayoutRequest) => {
    const approvals = payout.Approvals?.filter(a => a.Approved).length || 0
    const rejections = payout.Approvals?.filter(a => !a.Approved).length || 0

    if (payout.Status === 'approved') return { text: 'Approved', color: 'green' }
    if (payout.Status === 'rejected') return { text: 'Rejected', color: 'red' }
    if (rejections > 0) return { text: 'Rejected', color: 'red' }
    if (approvals >= 1) return { text: 'Approved', color: 'green' }
    return { text: `Pending (${approvals}/1 approval)`, color: 'yellow' }
  }

  const hasUserVoted = (payout: PayoutRequest) => {
    return payout.Approvals?.some(a => a.Admin.id === currentUser.id)
  }

  return (
    <div className="space-y-6">
      {/* Create Payout Request */}
      {isAdmin && group.Status === 'active' && (
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Payout Management</h3>
            <button
              onClick={() => setShowCreateModal(true)}
              className="btn btn-primary"
            >
              <DollarSign className="w-4 h-4 mr-2" />
              Create Payout Request
            </button>
          </div>
          
         