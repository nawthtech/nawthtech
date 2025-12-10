// Cloudflare Management Component
import React, { useState } from 'react';
import {
  Box, Card, Text, Group, Stack, Title, Button,
  Badge, Progress, Table, ScrollArea, ActionIcon,
  Modal, TextInput, Textarea, JsonInput, Alert,
  SimpleGrid, Switch, useMantineTheme
} from '@mantine/core';
import {
  IconCloud, IconDatabase, IconKey, IconSettings,
  IconRefresh, IconPlayerPlay, IconTrash,
  IconDownload, IconAlertCircle, IconCheck,
  IconX, IconServer, IconNetwork, IconShield
} from '@tabler/icons-react';
import { adminAPI } from '@/services/admin';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';

const CloudflareManagement: React.FC = () => {
  const theme = useMantineTheme();
  const queryClient = useQueryClient();
  const [purgeModal, setPurgeModal] = useState(false);
  const [deployModal, setDeployModal] = useState(false);
  const [workerScript, setWorkerScript] = useState('');

  // Fetch Cloudflare metrics
  const { data: cloudflareMetrics, isLoading } = useQuery({
    queryKey: ['cloudflare-metrics'],
    queryFn: () => adminAPI.getCloudflareMetrics().then(res => res.data),
    refetchInterval: 60000,
  });

  // Mutations
  const purgeCacheMutation = useMutation({
    mutationFn: () => adminAPI.purgeCloudflareCache(),
    onSuccess: () => {
      notifications.show({
        title: 'تم المسح',
        message: 'تم مسح ذاكرة التخزين المؤقت لـ Cloudflare',
        color: 'green',
      });
      queryClient.invalidateQueries({ queryKey: ['cloudflare-metrics'] });
      setPurgeModal(false);
    },
  });

  const deployWorkerMutation = useMutation({
    mutationFn: () => adminAPI.deployWorker('custom-worker', workerScript),
    onSuccess: () => {
      notifications.show({
        title: 'تم النشر',
        message: 'تم نشر Worker جديد بنجاح',
        color: 'green',
      });
      setDeployModal(false);
      setWorkerScript('');
    },
  });

  // Worker templates
  const workerTemplates = [
    {
      name: 'API Gateway',
      description: 'بوابة API مع التخزين المؤقت والمصادقة',
      script: `export default {
  async fetch(request, env) {
    // API Gateway implementation
    const url = new URL(request.url);
    
    // Add caching headers
    const headers = new Headers();
    headers.set('Cache-Control', 'public, max-age=3600');
    
    return new Response('Hello from Cloudflare Worker!', { headers });
  }
};`
    },
    {
      name: 'Image Optimizer',
      description: 'تحسين الصور تلقائياً',
      script: `export default {
  async fetch(request, env) {
    // Image optimization logic
    return fetch(request);
  }
};`
    }
  ];

  return (
    <Box p="md">
      <Group justify="space-between" mb="xl">
        <div>
          <Title order={2}>إدارة Cloudflare</Title>
          <Text c="dimmed" size="sm">مراقبة وإدارة خدمات Cloudflare</Text>
        </div>
        
        <Group>
          <Button
            leftSection={<IconRefresh size={16} />}
            variant="light"
            onClick={() => queryClient.invalidateQueries({ queryKey: ['cloudflare-metrics'] })}
          >
            تحديث
          </Button>
          <Button
            leftSection={<IconShield size={16} />}
            color="orange"
            onClick={() => setPurgeModal(true)}
          >
            مسح الكاش
          </Button>
          <Button
            leftSection={<IconCloud size={16} />}
            color="blue"
            onClick={() => setDeployModal(true)}
          >
            نشر Worker
          </Button>
        </Group>
      </Group>

      {/* Cloudflare Services Status */}
      <SimpleGrid cols={{ base: 1, sm: 2, lg: 4 }} spacing="md" mb="md">
        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>CDN</Text>
            <Badge color="green" variant="light">
              نشط
            </Badge>
          </Group>
          <Text size="xl" fw={700}>
            {cloudflareMetrics?.analytics?.requests || 0}
          </Text>
          <Text size="sm" c="dimmed">طلبات اليوم</Text>
          <Progress 
            value={cloudflareMetrics?.analytics?.cacheHitRate || 0} 
            mt="md" 
            size="sm" 
            color="blue" 
          />
          <Text size="xs" mt="xs">معدل ضربات الكاش: {cloudflareMetrics?.analytics?.cacheHitRate || 0}%</Text>
        </Card>

        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>Workers</Text>
            <Badge color="green" variant="light">
              {cloudflareMetrics?.workers?.active || 0} نشط
            </Badge>
          </Group>
          <Text size="xl" fw={700}>
            {cloudflareMetrics?.workers?.invocations || 0}
          </Text>
          <Text size="sm" c="dimmed">استدعاءات اليوم</Text>
          <Text size="xs" mt="xs">
            متوسط الاستجابة: {cloudflareMetrics?.workers?.avgResponseTime || 0}ms
          </Text>
        </Card>

        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>D1 Database</Text>
            <Badge color="green" variant="light">
              نشط
            </Badge>
          </Group>
          <Text size="xl" fw={700}>
            {cloudflareMetrics?.d1?.queries || 0}
          </Text>
          <Text size="sm" c="dimmed">استعلامات اليوم</Text>
          <Text size="xs" mt="xs">
            حجم البيانات: {(cloudflareMetrics?.d1?.size || 0).toFixed(2)} MB
          </Text>
        </Card>

        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>KV Storage</Text>
            <Badge color="green" variant="light">
              {cloudflareMetrics?.kv?.keys || 0} مفتاح
            </Badge>
          </Group>
          <Text size="xl" fw={700">
            {cloudflareMetrics?.kv?.operations || 0}
          </Text>
          <Text size="sm" c="dimmed">عمليات اليوم</Text>
          <Text size="xs" mt="xs">
            معدل القراءة: {cloudflareMetrics?.kv?.readRate || 0}/ثانية
          </Text>
        </Card>
      </SimpleGrid>

      {/* Workers Management */}
      <Card withBorder mb="md">
        <Group justify="space-between" mb="md">
          <Title order={4}>Cloudflare Workers</Title>
          <Button 
            size="xs" 
            leftSection={<IconCloud size={14} />}
            onClick={() => setDeployModal(true)}
          >
            نشر Worker جديد
          </Button>
        </Group>

        <Table>
          <thead>
            <tr>
              <th>الاسم</th>
              <th>الحالة</th>
              <th>الطلبات</th>
              <th>متوسط الاستجابة</th>
              <th>تاريخ النشر</th>
              <th>الإجراءات</th>
            </tr>
          </thead>
          <tbody>
            {cloudflareMetrics?.workers?.list?.map((worker: any) => (
              <tr key={worker.id}>
                <td>
                  <Group gap="xs">
                    <IconServer size={16} />
                    <Text>{worker.name}</Text>
                  </Group>
                </td>
                <td>
                  <Badge 
                    color={worker.status === 'active' ? 'green' : 'yellow'}
                    variant="light"
                  >
                    {worker.status === 'active' ? 'نشط' : 'متوقف'}
                  </Badge>
                </td>
                <td>{worker.requests}</td>
                <td>{worker.avgResponseTime}ms</td>
                <td>
                  <Text size="sm">
                    {new Date(worker.deployedAt).toLocaleDateString('ar-SA')}
                  </Text>
                </td>
                <td>
                  <Group gap="xs">
                    <ActionIcon variant="light" size="sm" color="blue">
                      <IconPlayerPlay size={14} />
                    </ActionIcon>
                    <ActionIcon variant="light" size="sm" color="yellow">
                      <IconSettings size={14} />
                    </ActionIcon>
                    <ActionIcon variant="light" size="sm" color="red">
                      <IconTrash size={14} />
                    </ActionIcon>
                  </Group>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      </Card>

      {/* Modals */}
      <Modal
        opened={purgeModal}
        onClose={() => setPurgeModal(false)}
        title="مسح ذاكرة التخزين المؤقت"
        size="md"
      >
        <Stack gap="md">
          <Alert color="orange" icon={<IconAlertCircle size={16} />}>
            <Text size="sm">
              سيتم مسح ذاكرة التخزين المؤقت لـ Cloudflare بالكامل.
              قد يؤثر هذا على الأداء مؤقتاً.
            </Text>
          </Alert>

          <Group>
            <Switch label="مسح كل المنطقة" defaultChecked />
            <Switch label="مسح ملفات معينة" />
          </Group>

          <Textarea
            label="ملفات محددة (اختياري)"
            placeholder="أدخل مسارات الملفات لمسحها..."
            autosize
            minRows={3}
          />

          <Group justify="flex-end">
            <Button variant="light" onClick={() => setPurgeModal(false)}>
              إلغاء
            </Button>
            <Button 
              color="red" 
              loading={purgeCacheMutation.isPending}
              onClick={() => purgeCacheMutation.mutate()}
            >
              تأكيد المسح
            </Button>
          </Group>
        </Stack>
      </Modal>

      <Modal
        opened={deployModal}
        onClose={() => setDeployModal(false)}
        title="نشر Cloudflare Worker"
        size="xl"
      >
        <Tabs defaultValue="templates">
          <Tabs.List>
            <Tabs.Tab value="templates">القوالب</Tabs.Tab>
            <Tabs.Tab value="custom">مخصص</Tabs.Tab>
          </Tabs.List>

          <Tabs.Panel value="templates" pt="md">
            <Stack gap="md">
              {workerTemplates.map((template) => (
                <Card key={template.name} withBorder>
                  <Group justify="space-between" mb="xs">
                    <div>
                      <Text fw={500}>{template.name}</Text>
                      <Text size="sm" c="dimmed">{template.description}</Text>
                    </div>
                    <Button
                      size="xs"
                      onClick={() => {
                        setWorkerScript(template.script);
                      }}
                    >
                      استخدام القالب
                    </Button>
                  </Group>
                  <ScrollArea style={{ height: 100 }}>
                    <Text size="xs" c="dimmed" style={{ fontFamily: 'monospace' }}>
                      {template.script}
                    </Text>
                  </ScrollArea>
                </Card>
              ))}
            </Stack>
          </Tabs.Panel>

          <Tabs.Panel value="custom" pt="md">
            <Stack gap="md">
              <TextInput
                label="اسم Worker"
                placeholder="my-custom-worker"
                required
              />
              
              <JsonInput
                label="متغيرات البيئة"
                placeholder='{ "API_KEY": "your-key" }'
                formatOnBlur
                autosize
                minRows={3}
              />
              
              <Textarea
                label="كود Worker"
                placeholder="أدخل كود JavaScript لـ Worker..."
                value={workerScript}
                onChange={(e) => setWorkerScript(e.target.value)}
                autosize
                minRows={10}
                styles={{
                  input: {
                    fontFamily: 'monospace',
                    fontSize: 12,
                  }
                }}
              />

              <Alert color="blue">
                <Text size="sm">
                  يمكنك استخدام المتغيرات البيئية عبر <code>env</code> في الكود
                </Text>
              </Alert>
            </Stack>
          </Tabs.Panel>
        </Tabs>

        <Group justify="flex-end" mt="xl">
          <Button variant="light" onClick={() => setDeployModal(false)}>
            إلغاء
          </Button>
          <Button 
            color="blue" 
            loading={deployWorkerMutation.isPending}
            onClick={() => deployWorkerMutation.mutate()}
            disabled={!workerScript.trim()}
          >
            نشر Worker
          </Button>
        </Group>
      </Modal>
    </Box>
  );
};

export default CloudflareManagement;