import { useEffect, useState } from 'react'
import { Copy, KeyRound, Plus, Trash2 } from 'lucide-react'
import { AppService } from '@/services/app-service'
import type { App, AppAPIKey, AppSecret } from '@/types/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogScrollBody, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { CreateAppSecretDialog } from '@/features/config/components/create-app-secret-dialog'
import { scopeLabel } from '@/features/config/components/app-secret-scopes'
import { CreateApiKeyDialog } from '@/features/config/components/create-api-key-dialog'
import { DeleteApiKeyDialog } from '@/features/config/components/delete-api-key-dialog'
import { toast } from 'sonner'

interface APIKeysDialogProps { app: App; open: boolean; onOpenChange: (open: boolean) => void }

export function APIKeysDialog({ app, open, onOpenChange }: APIKeysDialogProps) {
  const [secrets, setSecrets] = useState<AppSecret[]>([])
  const [loading, setLoading] = useState(false)
  const [createOpen, setCreateOpen] = useState(false)
  const [apiKeys, setApiKeys] = useState<AppAPIKey[]>([])
  const [createApiKeyOpen, setCreateApiKeyOpen] = useState(false)
  const [deleteApiKeyOpen, setDeleteApiKeyOpen] = useState(false)
  const [selectedApiKey, setSelectedApiKey] = useState<AppAPIKey | null>(null)
  const load = async () => {
    try {
      setLoading(true)
      const [nextSecrets, nextApiKeys] = await Promise.all([
        app.role === 'owner' ? AppService.getAppSecrets(app.id) : Promise.resolve([]),
        AppService.getAppAPIKeys(app.id),
      ])
      setSecrets(nextSecrets)
      setApiKeys(nextApiKeys)
    }
    catch (error) { toast.error((error as Error).message || '加载应用凭证失败') }
    finally { setLoading(false) }
  }
  useEffect(() => { if (open) load() }, [open, app.id]) // eslint-disable-line react-hooks/exhaustive-deps

  return <>
    <Dialog open={open} onOpenChange={onOpenChange}><DialogContent className="max-h-[80vh] overflow-hidden sm:max-w-[760px]"><DialogHeader><DialogTitle className="flex items-center gap-2"><KeyRound className="h-5 w-5" />应用凭证</DialogTitle><DialogDescription>App Key 用于客户端 SDK；App Secret 仅用于服务端 API。</DialogDescription></DialogHeader><DialogScrollBody className="space-y-5">
      <div className="space-y-2"><div className="text-sm font-medium">App Key</div><div className="flex gap-2"><Input value={app.app_key || ''} readOnly className="font-mono" /><Button size="icon" variant="outline" title="复制 App Key" onClick={() => navigator.clipboard.writeText(app.app_key)}><Copy className="h-4 w-4" /></Button></div></div>
      {app.role === 'owner' && <div className="space-y-3"><div className="flex items-center justify-between"><div className="font-medium">App Secrets</div><Button size="sm" onClick={() => setCreateOpen(true)}><Plus className="mr-2 h-4 w-4" />创建</Button></div>{loading ? <div className="py-6 text-center text-muted-foreground">加载中...</div> : secrets.length === 0 ? <div className="py-6 text-center text-muted-foreground">尚未创建 App Secret</div> : <Table><TableHeader><TableRow><TableHead>名称</TableHead><TableHead>权限</TableHead><TableHead>状态</TableHead></TableRow></TableHeader><TableBody>{secrets.map((secret) => <TableRow key={secret.id}><TableCell><div>{secret.name}</div><code className="text-xs text-muted-foreground">{secret.prefix}...{secret.suffix}</code></TableCell><TableCell><div className="flex flex-wrap gap-1">{secret.scopes.map((scope) => <Badge key={scope} variant="secondary">{scopeLabel(scope)}</Badge>)}</div></TableCell><TableCell>{secret.status === 1 && !secret.revoked_at ? '有效' : '已撤销'}</TableCell></TableRow>)}</TableBody></Table>}</div>}
      <div className="space-y-3"><div className="flex items-start justify-between gap-4"><div><div className="font-medium">API 密钥</div><div className="text-xs text-muted-foreground">保留管理能力，当前不参与认证。</div></div><Button size="sm" variant="outline" onClick={() => setCreateApiKeyOpen(true)}><Plus className="mr-2 h-4 w-4" />创建</Button></div>{loading ? <div className="py-6 text-center text-muted-foreground">加载中...</div> : apiKeys.length === 0 ? <div className="py-6 text-center text-muted-foreground">暂无 API 密钥</div> : <Table><TableHeader><TableRow><TableHead>名称</TableHead><TableHead>密钥</TableHead><TableHead className="text-right">操作</TableHead></TableRow></TableHeader><TableBody>{apiKeys.map((apiKey) => <TableRow key={apiKey.id}><TableCell>{apiKey.name}</TableCell><TableCell><code className="text-xs text-muted-foreground">{apiKey.key_prefix}...{apiKey.key_suffix}</code></TableCell><TableCell className="text-right"><Button size="icon" variant="ghost" title="删除 API 密钥" onClick={() => { setSelectedApiKey(apiKey); setDeleteApiKeyOpen(true) }}><Trash2 className="h-4 w-4" /></Button></TableCell></TableRow>)}</TableBody></Table>}</div>
    </DialogScrollBody><Button onClick={() => onOpenChange(false)}>关闭</Button></DialogContent></Dialog>
    <CreateAppSecretDialog app={app} open={createOpen} onOpenChange={setCreateOpen} onSuccess={load} />
    <CreateApiKeyDialog app={app} open={createApiKeyOpen} onOpenChange={setCreateApiKeyOpen} onSuccess={load} />
    <DeleteApiKeyDialog app={app} apiKey={selectedApiKey} open={deleteApiKeyOpen} onOpenChange={setDeleteApiKeyOpen} onSuccess={load} />
  </>
}
