import React, { useState } from 'react'
import { Bell, Check, X, Users, DollarSign, UserPlus } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api, { notificationApi } from '../services/api'
import type { Notification, GroupInvitation } from '../types'
import toast from 'react-hot-toast'

interface NotificationCenterProps {
  isCollapsed: boolean;
}
