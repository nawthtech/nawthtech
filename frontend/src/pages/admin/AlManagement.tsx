import React, { useState } from 'react';
import {
  Box, Card, Text, Group, Stack, Table, ScrollArea,
  Title, Button, Badge, Progress, Select, TextInput,
  NumberInput, Switch, Tabs, ActionIcon, Modal,
  Textarea, JsonInput, LoadingOverlay, Alert,
  SimpleGrid, RingProgress, useMantineTheme
} from '@mantine/core';
import {
  IconRobot, IconVideo, IconImage, IconBrain,
  IconSettings, IconRefresh, IconPlayerPlay,
  IconPlayerStop, IconTrash, IconDownload,
  IconChartBar, IconDatabase, IconCloud,
  IconAlertCircle, IconCheck, IconX
} from '@tabler/icons-react';
import { adminAPI, adminHelpers } from '@/services/admin';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';

const AIManagement: React.FC = () => {
  const theme = useMantineTheme();
  const queryClient = useQueryClient();
  const [selectedModel, setSelectedModel] = useState('gemini');
  const [testModal, setTestModal] = useState(false);
  const [testPrompt, setTestPrompt] = useState('');
  const [configModal, setConfigModal] = useState(false);
  const [modelConfig, setModelConfig] = useState<any>({});

  // Fetch AI metrics
  const { data: aiMetrics, isLoading } = useQuery({
    queryKey: ['ai-metrics'],
    queryFn: () => adminAPI.getAIDashboardMetrics().then(res => res.data),
  });

  // Fetch AI requests
  const { data: aiRequests } = useQuery({
    queryKey: ['ai-requests'],
    queryFn: () => adminAPI.getAIRequests({ limit: 50 }).then(res => res.data),
  });

  // Mutations
  const retryAIMutation = useMutation({
    mutationFn: (requestId: string) => adminAPI.retryAIRequest(requestId),
    onSuccess: () => {
      notifications.show({
        title: 'تمت إعادة المحاولة',
        message: 'تم إعادة معالجة الطلب',
        color: 'green',
      });
      queryClient.invalidateQueries({ queryKey: ['ai-requests'] });
    },
  });

  const cancelAIMutation = useMutation({
    mutationFn: (requestId: string) => adminAPI.cancelAIRequest(requestId),
    onSuccess: () => {
      notifications.show({
        title: 'تم الإلغاء',
        message: 'تم إلغاء طلب الذكاء الاصطناعي',
        color: 'yellow',
      });
      queryClient.invalidateQueries({ queryKey: ['ai-requests'] });
    },
  });

  const updateConfigMutation = useMutation({
    mutationFn: () => adminAPI.updateAISettings(modelConfig),
    onSuccess: () => {
      notifications.show({
        title: 'تم التحديث',
        message: 'تم تحديث إعدادات الذكاء الاصطناعي',
        color: 'green',
      });
      setConfigModal(false);
      queryClient.invalidateQueries({ queryKey: ['ai-metrics'] });
    },
  });

  // Model costs
  const modelCosts = [
    { model: 'Gemini', cost: '$0.001/request', type: 'Text/Image' },
    { model: 'OpenAI GPT-4', cost: '$0.002/request', type: 'Text' },
    { model: 'Stability AI', cost: '$0.01/image', type: 'Image' },
    { model: 'Luma AI', cost: '$0.02/video', type: 'Video' },
    { model: 'Runway ML', cost: '$0.03/video', type: 'Video' },
    { model: 'Pika Labs', cost: '$0.015/video', type: 'Video' },
  ];

  return (
    <Box p="md">
      <Group justify="space-between" mb="xl">
        <div>
          <Title order={2}>إدارة الذكاء الاصطناعي</Title>
          <Text c="dimmed" size="sm">مراقبة وإدارة نماذج الذكاء الاصطناعي</Text>
        </div>
        
        <Group>
          <Button
            leftSection={<IconSettings size={16} />}
            onClick={() => setConfigModal(true)}
            variant="light"
          >
            الإعدادات
          </Button>
          <Button
            leftSection={<IconRobot size={16} />}
            onClick={() => setTestModal(true)}
            color="grape"
          >
            اختبار النماذج
          </Button>
        </Group>
      </Group>

      {/* AI Metrics Cards */}
      <SimpleGrid cols={{ base: 1, sm: 2, lg: 4 }} spacing="md" mb="md">
        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>إجمالي الطلبات</Text>
            <Badge color="grape" variant="light">
              AI
            </Badge>
          </Group>
          <Text size="xl" fw={700}>
            {adminHelpers.formatNumber(aiMetrics?.totalRequests || 0)}
          </Text>
          <Progress 
            value={(aiMetrics?.totalSuccessfulRequests || 0) / (aiMetrics?.totalRequests || 1) * 100} 
            mt="md" 
            size="sm" 
            color="grape" 
          />
        </Card>

        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>معدل النجاح</Text>
            <Badge color="green" variant="light">
              {((aiMetrics?.totalSuccessfulRequests || 0) / (aiMetrics?.totalRequests || 1) * 100).toFixed(1)}%
            </Badge>
          </Group>
          <RingProgress
            size={80}
            thickness={8}
            sections={[
              { 
                value: (aiMetrics?.totalSuccessfulRequests || 0) / (aiMetrics?.totalRequests || 1) * 100, 
                color: 'green' 
              }
            ]}
            label={
              <Text fw={700} size="xl" ta="center">
                {((aiMetrics?.totalSuccessfulRequests || 0) / (aiMetrics?.totalRequests || 1) * 100).toFixed(1)}%
              </Text>
            }
          />
        </Card>

        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>متوسط وقت المعالجة</Text>
            <IconBrain size={20} color={theme.colors.blue[6]} />
          </Group>
          <Text size="xl" fw={700}>
            {aiMetrics?.averageProcessingTime?.toFixed(2) || 0} ثانية
          </Text>
          <Text size="sm" c="dimmed" mt="xs">
            أسرع نموذج: {Math.min(
              aiMetrics?.gemini?.averageTime || 0,
              aiMetrics?.openai?.averageTime || 0,
              aiMetrics?.ollama?.averageTime || 0
            ).toFixed(2)} ثانية
          </Text>
        </Card>

        <Card withBorder>
          <Group justify="space-between" mb="md">
            <Text fw={500}>التكلفة الإجمالية</Text>
            <Badge color="red" variant="light">
              {adminHelpers.formatCurrency(aiMetrics?.totalCost || 0, 'USD')}
            </Badge>
          </Group>
          <Text size="xl" fw={700}>
            {adminHelpers.formatCurrency(aiMetrics?.totalCost || 0, 'USD')}
          </Text>
          <Text size="sm" c="dimmed" mt="xs">
            اليوم: {adminHelpers.formatCurrency(aiMetrics?.totalCost ? aiMetrics.totalCost / 30 : 0, 'USD')}
          </Text>
        </Card>
      </SimpleGrid>

      {/* Models Overview */}
      <Tabs defaultValue="models" mb="md">
        <Tabs.List>
          <Tabs.Tab value="models" leftSection={<IconRobot size={16} />}>
            النماذج
          </Tabs.Tab>
          <Tabs.Tab value="requests" leftSection={<IconDatabase size={16} />}>
            الطلبات الأخيرة
          </Tabs.Tab>
          <Tabs.Tab value="video" leftSection={<IconVideo size={16} />}>
            توليد الفيديو
          </Tabs.Tab>
          <Tabs.Tab value="costs" leftSection={<IconChartBar size={16} />}>
            التكاليف
          </Tabs.Tab>
        </Tabs.List>

        <Tabs.Panel value="models" pt="md">
          <Card withBorder>
            <Table>
              <thead>
                <tr>
                  <th>النموذج</th>
                  <th>النوع</th>
                  <th>طلبات اليوم</th>
                  <th>معدل النجاح</th>
                  <th>متوسط الوقت</th>
                  <th>الحالة</th>
                  <th>الإجراءات</th>
                </tr>
              </thead>
              <tbody>
                {aiMetrics && Object.entries(aiMetrics).map(([key, model]: [string, any]) => (
                  typeof model === 'object' && model.requests !== undefined && (
                    <tr key={key}>
                      <td>
                        <Group gap="xs">
                          <Badge 
                            variant="light" 
                            color={
                              key === 'gemini' ? 'grape' :
                              key === 'openai' ? 'blue' :
                              key === 'stability' ? 'pink' :
                              key === 'videoGeneration' ? 'cyan' : 'gray'
                            }
                          >
                            {adminHelpers.formatAIModelName(key)}
                          </Badge>
                        </Group>
                      </td>
                      <td>
                        <Text size="sm">
                          {key === 'videoGeneration' ? 'فيديو' :
                           key === 'stability' ? 'صورة' : 'نص'}
                        </Text>
                      </td>
                      <td>{model.requests || 0}</td>
                      <td>
                        <Group>
                          <Text>{model.successRate?.toFixed(1) || 0}%</Text>
                          <Progress 
                            value={model.successRate || 0} 
                            size="sm" 
                            w={100} 
                            color={
                              model.successRate > 90 ? 'green' :
                              model.successRate > 70 ? 'yellow' : 'red'
                            }
                          />
                        </Group>
                      </td>
                      <td>{model.averageTime?.toFixed(2) || 0} ثانية</td>
                      <td>
                        <Badge 
                          color={model.successRate > 90 ? 'green' : model.successRate > 70 ? 'yellow' : 'red'}
                          variant="light"
                        >
                          {model.successRate > 90 ? 'ممتاز' : model.successRate > 70 ? 'جيد' : 'تحتاج تحسين'}
                        </Badge>
                      </td>
                      <td>
                        <Group gap="xs">
                          <ActionIcon 
                            variant="light" 
                            size="sm" 
                            color="blue"
                            onClick={() => {
                              setSelectedModel(key);
                              setTestModal(true);
                            }}
                          >
                            <IconPlayerPlay size={14} />
                          </ActionIcon>
                          <ActionIcon variant="light" size="sm" color="yellow">
                            <IconRefresh size={14} />
                          </ActionIcon>
                        </Group>
                      </td>
                    </tr>
                  )
                ))}
              </tbody>
            </Table>
          </Card>
        </Tabs.Panel>

        <Tabs.Panel value="requests" pt="md">
          <Card withBorder>
            <ScrollArea style={{ height: 500 }}>
              <Table>
                <thead>
                  <tr>
                    <th>المعرف</th>
                    <th>النموذج</th>
                    <th>النوع</th>
                    <th>المستخدم</th>
                    <th>الحالة</th>
                    <th>الوقت</th>
                    <th>الإجراءات</th>
                  </tr>
                </thead>
                <tbody>
                  {aiRequests?.map((request: any) => (
                    <tr key={request.id}>
                      <td>
                        <Text size="sm" truncate style={{ maxWidth: 100 }}>
                          {request.id}
                        </Text>
                      </td>
                      <td>
                        <Badge variant="light" color="grape" size="sm">
                          {request.model}
                        </Badge>
                      </td>
                      <td>
                        <Badge variant="light" color="blue" size="sm">
                          {request.type}
                        </Badge>
                      </td>
                      <td>
                        <Text size="sm">{request.user}</Text>
                      </td>
                      <td>
                        <Badge 
                          color={
                            request.status === 'completed' ? 'green' :
                            request.status === 'processing' ? 'blue' :
                            request.status === 'failed' ? 'red' : 'yellow'
                          }
                          variant="light"
                        >
                          {request.status === 'completed' ? 'مكتمل' :
                           request.status === 'processing' ? 'جاري' :
                           request.status === 'failed' ? 'فشل' : 'معلق'}
                        </Badge>
                      </td>
                      <td>
                        <Text size="sm">
                          {new Date(request.timestamp).toLocaleTimeString('ar-SA')}
                        </Text>
                      </td>
                      <td>
                        <Group gap="xs">
                          {request.status === 'failed' && (
                            <ActionIcon 
                              variant="light" 
                              size="sm" 
                              color="green"
                              onClick={() => retryAIMutation.mutate(request.id)}
                            >
                              <IconRefresh size={14} />
                            </ActionIcon>
                          )}
                          {(request.status === 'processing' || request.status === 'pending') && (
                            <ActionIcon 
                              variant="light" 
                              size="sm" 
                              color="red"
                              onClick={() => cancelAIMutation.mutate(request.id)}
                            >
                              <IconPlayerStop size={14} />
                            </ActionIcon>
                          )}
                          <ActionIcon variant="light" size="sm" color="blue">
                            <IconEye size={14} />
                          </ActionIcon>
                        </Group>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </Table>
            </ScrollArea>
          </Card>
        </Tabs.Panel>

        <Tabs.Panel value="costs" pt="md">
          <SimpleGrid cols={{ base: 1, md: 2 }} spacing="md">
            <Card withBorder>
              <Title order={4} mb="md">تكاليف النماذج</Title>
              <Table>
                <thead>
                  <tr>
                    <th>النموذج</th>
                    <th>النوع</th>
                    <th>التكلفة/طلب</th>
                    <th>التكلفة الشهرية</th>
                  </tr>
                </thead>
                <tbody>
                  {modelCosts.map((model) => (
                    <tr key={model.model}>
                      <td>{model.model}</td>
                      <td>
                        <Badge variant="light" size="sm">
                          {model.type}
                        </Badge>
                      </td>
                      <td>
                        <Text fw={500}>{model.cost}</Text>
                      </td>
                      <td>
                        <Text fw={500}>
                          {adminHelpers.formatCurrency(
                            adminHelpers.estimateAICost(
                              model.model.toLowerCase().split(' ')[0],
                              1000,
                              model.type === 'Video' ? 'video' : 
                              model.type === 'Image' ? 'image' : 'text'
                            ),
                            'USD'
                          )}
                        </Text>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </Table>
            </Card>

            <Card withBorder>
              <Title order={4} mb="md">التحكم في التكاليف</Title>
              <Stack gap="md">
                <NumberInput
                  label="حد الطلبات اليومي"
                  description="الحد الأقصى لطلبات الذكاء الاصطناعي اليومية"
                  min={0}
                  defaultValue={1000}
                />
                <NumberInput
                  label="الحد الشهري ($)"
                  description="الحد الأقصى للتكلفة الشهرية"
                  min={0}
                  prefix="$"
                  defaultValue={500}
                />
                <Switch
                  label="إيقاف التلقائي عند الوصول للحد"
                  defaultChecked
                />
                <Switch
                  label="إرسال تنبيهات عند ارتفاع التكاليف"
                  defaultChecked
                />
                <Button color="red" variant="light">
                  تطبيق حدود التكلفة
                </Button>
              </Stack>
            </Card>
          </SimpleGrid>
        </Tabs.Panel>
      </Tabs>

      {/* Modals */}
      <Modal
        opened={testModal}
        onClose={() => setTestModal(false)}
        title="اختبار نموذج الذكاء الاصطناعي"
        size="lg"
      >
        <Stack gap="md">
          <Select
            label="اختر النموذج"
            data={[
              { label: 'Google Gemini', value: 'gemini' },
              { label: 'OpenAI GPT-4', value: 'openai' },
              { label: 'Stability AI (صور)', value: 'stability' },
              { label: 'Luma AI (فيديو)', value: 'luma' },
              { label: 'Runway ML (فيديو)', value: 'runway' },
              { label: 'Pika Labs (فيديو)', value: 'pika' },
            ]}
            value={selectedModel}
            onChange={(value) => value && setSelectedModel(value)}
          />

          <Textarea
            label="النص المطلوب"
            placeholder="أدخل النص المطلوب معالجته..."
            value={testPrompt}
            onChange={(e) => setTestPrompt(e.target.value)}
            autosize
            minRows={4}
          />

          {selectedModel.includes('video') && (
            <Stack gap="xs">
              <Text size="sm" fw={500}>إعدادات الفيديو</Text>
              <Group grow>
                <Select
                  label="المدة"
                  data={[
                    { label: '5 ثواني', value: '5' },
                    { label: '10 ثواني', value: '10' },
                    { label: '15 ثواني', value: '15' },
                    { label: '30 ثواني', value: '30' },
                  ]}
                  defaultValue="10"
                />
                <Select
                  label="الدقة"
                  data={[
                    { label: '480p', value: '480' },
                    { label: '720p', value: '720' },
                    { label: '1080p', value: '1080' },
                    { label: '4K', value: '2160' },
                  ]}
                  defaultValue="1080"
                />
              </Group>
            </Stack>
          )}

          <Alert color="blue" icon={<IconAlertCircle size={16} />}>
            <Text size="sm">
              التكلفة المتوقعة:{' '}
              {adminHelpers.estimateAICost(
                selectedModel,
                1,
                selectedModel.includes('video') ? 'video' : 
                selectedModel.includes('image') || selectedModel === 'stability' ? 'image' : 'text'
              )} دولار
            </Text>
          </Alert>

          <Button
            fullWidth
            color="grape"
            leftSection={<IconRobot size={16} />}
            onClick={() => {
              notifications.show({
                title: 'جاري المعالجة',
                message: 'جاري معالجة طلب الذكاء الاصطناعي',
                color: 'blue',
              });
              setTestModal(false);
            }}
          >
            بدء الاختبار
          </Button>
        </Stack>
      </Modal>

      <Modal
        opened={configModal}
        onClose={() => setConfigModal(false)}
        title="إعدادات الذكاء الاصطناعي"
        size="xl"
      >
        <Tabs defaultValue="general">
          <Tabs.List>
            <Tabs.Tab value="general">عام</Tabs.Tab>
            <Tabs.Tab value="models">النماذج</Tabs.Tab>
            <Tabs.Tab value="video">الفيديو</Tabs.Tab>
            <Tabs.Tab value="limits">الحدود</Tabs.Tab>
          </Tabs.List>

          <Tabs.Panel value="general" pt="md">
            <Stack gap="md">
              <Switch
                label="تمكين الذكاء الاصطناعي"
                defaultChecked
                onChange={(e) => setModelConfig({ ...modelConfig, enabled: e.target.checked })}
              />
              <Select
                label="النموذج الافتراضي"
                data={[
                  { label: 'Google Gemini', value: 'gemini' },
                  { label: 'OpenAI GPT-4', value: 'openai' },
                  { label: 'Ollama LLM', value: 'ollama' },
                ]}
                defaultValue="gemini"
                onChange={(value) => setModelConfig({ ...modelConfig, defaultModel: value })}
              />
              <NumberInput
                label="حد الطلبات في الثانية"
                description="الحد الأقصى لطلبات API في الثانية"
                min={1}
                max={100}
                defaultValue={10}
                onChange={(value) => setModelConfig({ ...modelConfig, rateLimit: value })}
              />
            </Stack>
          </Tabs.Panel>

          <Tabs.Panel value="models" pt="md">
            <Stack gap="md">
              <Text size="sm" fw={500}>تفعيل النماذج</Text>
              <Switch label="Google Gemini" defaultChecked />
              <Switch label="OpenAI GPT-4" defaultChecked />
              <Switch label="Ollama LLM" defaultChecked />
              <Switch label="Hugging Face Models" defaultChecked />
              <Switch label="Stability AI" defaultChecked />

              <Text size="sm" fw={500} mt="md">مفاتيح API</Text>
              <TextInput label="مفتاح Gemini API" type="password" />
              <TextInput label="مفتاح OpenAI API" type="password" />
              <TextInput label="مفتاح Stability API" type="password" />
            </Stack>
          </Tabs.Panel>

          <Tabs.Panel value="video" pt="md">
            <Stack gap="md">
              <Switch label="تمكين توليد الفيديو" defaultChecked />
              <Select
                label="منصة الفيديو الافتراضية"
                data={[
                  { label: 'Luma AI', value: 'luma' },
                  { label: 'Runway ML', value: 'runway' },
                  { label: 'Pika Labs', value: 'pika' },
                  { label: 'Gemini Veo', value: 'gemini-veo' },
                ]}
                defaultValue="luma"
              />
              <NumberInput
                label="الحد الأقصى لمدة الفيديو (ثواني)"
                min={5}
                max={60}
                defaultValue={30}
              />
              <Select
                label="الدقة الافتراضية"
                data={[
                  { label: '720p', value: '720' },
                  { label: '1080p', value: '1080' },
                  { label: '4K', value: '2160' },
                ]}
                defaultValue="1080"
              />
            </Stack>
          </Tabs.Panel>
        </Tabs>

        <Group justify="flex-end" mt="xl">
          <Button variant="light" onClick={() => setConfigModal(false)}>
            إلغاء
          </Button>
          <Button 
            color="blue" 
            onClick={() => updateConfigMutation.mutate()}
            loading={updateConfigMutation.isPending}
          >
            حفظ التغييرات
          </Button>
        </Group>
      </Modal>
    </Box>
  );
};

export default AIManagement;