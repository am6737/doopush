import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2, Key } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogScrollBody,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useAuthStore } from '@/stores/auth-store'
import { AppService } from '@/services/app-service'
import { requireApp } from '@/utils/app-utils'
import { toast } from 'sonner'
import { CreatedCredentialDialog } from './created-credential-dialog'
import type { App } from '@/types/api'

// 表单验证规则
const createApiKeySchema = z.object({
  name: z.string().min(1, '请输入密钥名称').max(100, '密钥名称不能超过100个字符'),
})

type CreateApiKeyFormData = z.infer<typeof createApiKeySchema>

interface CreateApiKeyDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess: () => void
  app?: App | null
}

export function CreateApiKeyDialog({ open, onOpenChange, onSuccess, app }: CreateApiKeyDialogProps) {
  const { currentApp } = useAuthStore()
  const targetApp = app ?? currentApp
  const [loading, setLoading] = useState(false)
  const [createdKey, setCreatedKey] = useState<{ api_key: string; warning?: string } | null>(null)
  const [showSuccessDialog, setShowSuccessDialog] = useState(false)

  const form = useForm<CreateApiKeyFormData>({
    resolver: zodResolver(createApiKeySchema),
    defaultValues: {
      name: '',
    },
  })

  const onSubmit = async (data: CreateApiKeyFormData) => {
    if (!requireApp(targetApp)) {
      return
    }

    try {
      setLoading(true)
      const result = await AppService.createAPIKey(targetApp.id, data)
      
      // 设置创建成功的密钥信息
      setCreatedKey({
        api_key: result.api_key,
        warning: result.warning
      })
      
      // 关闭创建对话框，显示成功对话框
      onOpenChange(false)
      setShowSuccessDialog(true)
      form.reset()
      onSuccess()
    } catch (error) {
      toast.error((error as Error).message || '创建API密钥失败')
    } finally {
      setLoading(false)
    }
  }

  const handleClose = () => {
    if (!loading) {
      form.reset()
      onOpenChange(false)
    }
  }

  const handleSuccessClose = () => {
    setShowSuccessDialog(false)
    setCreatedKey(null)
  }

  return (
    <>
      {/* 创建API密钥对话框 */}
      <Dialog open={open} onOpenChange={handleClose}>
        <DialogContent className="sm:max-w-[425px] max-h-[90vh] overflow-hidden flex flex-col">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Key className="h-5 w-5" />
              创建API密钥
            </DialogTitle>
            <DialogDescription>
              创建并保留 API 密钥。完整密钥只展示一次，请妥善保存。
            </DialogDescription>
          </DialogHeader>

          <DialogScrollBody>
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>密钥名称 *</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="例如：iOS生产环境密钥"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </form>
            </Form>
          </DialogScrollBody>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={handleClose} disabled={loading}>
              取消
            </Button>
            <Button onClick={form.handleSubmit(onSubmit)} disabled={loading}>
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              创建密钥
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <CreatedCredentialDialog
        open={showSuccessDialog}
        title="API 密钥创建成功"
        description="请立即复制并保存您的 API 密钥，此密钥不会再次显示。"
        credentialLabel="您的 API 密钥"
        credential={createdKey?.api_key ?? null}
        warning={createdKey?.warning}
        onClose={handleSuccessClose}
      />
    </>
  )
}
