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
   