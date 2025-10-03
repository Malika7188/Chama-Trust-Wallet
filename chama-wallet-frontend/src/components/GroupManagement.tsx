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
