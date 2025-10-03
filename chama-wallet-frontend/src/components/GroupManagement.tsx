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
  