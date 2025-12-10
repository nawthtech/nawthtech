import React, { useState, useEffect } from 'react';
import { 
  Box, Grid, Card, Text, Group, Stack, LoadingOverlay, 
  Title, Select, DatePicker, SegmentedControl, Button,
  Badge, Progress, RingProgress, Paper, SimpleGrid,
  useMantineTheme, Alert, Tabs, Table, ScrollArea,
  Menu, ActionIcon, Modal, TextInput, Textarea,
  NumberInput, Switch, MultiSelect, Timeline
} from '@mantine/core';
import {
  IconDashboard,
  IconUsers,
  IconShoppingCart,
  IconRobot,
  IconVideo,
  IconChartBar,
  IconSettings,
  IconAlertCircle,
  IconCloud,
  IconDatabase,
  IconRefresh,
  IconDownload,
  IconFilter,
  IconTrendingUp,
  IconTrendingDown,
  IconCircleCheck,
  IconCircleX,
  IconExternalLink,
  IconEye,
  IconEdit,
  IconTrash,
  IconPlayerPlay,
  IconPlayerStop,
  IconMessageCircle,
  IconBellRinging,
  IconCloudUpload,
  IconServer
} from '@tabler/icons-react';
import { LineChart, Line, BarChart, Bar, PieChart, Pie, Cell, 
  XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
  RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar
} from 'recharts';
import { adminAPI, adminHelpers, type DashboardData, type AnalyticsFilters } from '@/services/admin';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { useServices } from '@/contexts/ServicesContext';

const AdminDashboard: React.FC = () => {
  const theme = useMantineTheme();
  const queryClient = useQueryClient();
  const { services } = useServices();
  const [filters, setFilters] = useState<AnalyticsFilters>({ timeRange: 'month' });
  const [selectedTab, setSelectedTab] = useState<string>('overview');
  const [exportModal, setExportModal] = useState(false);
  const [testAIModal, setTestAIModal] = useState(false);
  const [testPrompt, setTestPrompt] = useState('');
  const [selectedModel, setSelectedModel] = useState('gemini');

  // Fetch dashboard data
  const { data: dashboardData, isLoading: dashboardLoading, refetch } = useQuery({
    queryKey: ['admin-dashboard', filters],
    queryFn: () => adminAPI.getDashboardData(filters).then(res => res.data),
  });

  // Fetch system health
  const { data: systemHealth } = useQuery({
    queryKey: ['system-health'],
    queryFn: () => adminAPI.checkSystemHealth().then(res => res.data),
    refetchInterval: 30000,
  });

  // Mutations
  const clearCacheMutation = useMutation({
    mutationFn: () => adminAPI.clearCache(),
    onSuccess: () => {
      notifications.show({
        title: 'تم مسح الذاكرة',
        message: 'تم مسح ذاكرة التخزين المؤقت بنجاح',
        color: 'green',
      });
      queryClient.invalidateQueries({ queryKey: ['admin-dashboard'] });
    },
  });

  const testAIMutation = useMutation({
    mutationFn: () => adminAPI.testAIService(selectedModel, testPrompt),
    onSuccess: () => {
      notifications.show({
        title: 'اختبار الذكاء الاصطناعي',
        message: 'تم تنفيذ الاختبار بنجاح',
        color: 'green',
      });
      setTestAIModal(false);
    },
  });

  // Chart data
  const revenueChartData = adminHelpers.formatChartData(dashboardData || {} as DashboardData, 'line');
  const aiUsageChartData = adminHelpers.formatChartData(dashboardData || {} as DashboardData, 'pie');
  const performanceChartData = adminHelpers.formatChartData(dashboardData || {} as DashboardData, 'radar');

  // Service status colors
  const getServiceStatusColor = (service: any) => {
    if (!service?.status) return 'gray';
    
    switch (service.status) {
      case 'healthy': return 'green';
      case 'degraded': return 'yellow';
      case 'unhealthy': return 'red';
      default: return 'gray';
    }
  };

  // Format metrics card
  const renderMetricCard = (title: string, value: any, change?: number, icon?: React.ReactNode) => (
    <Card withBorder p="md" radius="md">
      <Group justify="space-between" mb="xs">
        <Text size="sm" c="dimmed">{title}</Text>
        {icon}
      </Group>
      <Group align="flex-end" gap="xs">
        <Text size="xl" fw={700}>
          {typeof value === 'number' ? adminHelpers.formatNumber(value) : value}
        </Text>
        {change !== undefined && (
          <Badge 
            color={change > 0 ? 'green' : change < 0 ? 'red' : 'gray'}
            variant="light"
            leftSection={change > 0 ? <IconTrendingUp size={12} /> : <IconTrendingDown size={12} />}
          >
            {Math.abs(change)}%
          </Badge>
        )}
      </Group>
    </Card>
  );

  if (dashboardLoading) {
    return (
      <Box style={{ height: '80vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <LoadingOverlay visible />
      </Box>
    );
  }

  return (
    <Box p="md">
      {/* Header */}
      <Group justify="space-between" mb="xl">
        <div>
          <Title order={2}>لوحة تحكم المدير</Title>
          <Text c="dimmed" size="sm">مراقبة وتحليل أداء المنصة</Text>
        </div>
        
        <Group>
          <Button
            leftSection={<IconRefresh size={16} />}
            onClick={() => refetch()}
            variant="light"
          >
            تحديث البيانات
          </Button>
          
          <Button
            leftSection={<IconDownload size={16} />}
            onClick={() => setExportModal(true)}
            variant="filled"
            color="blue"
          >
            تصدير تقرير
          </Button>
          
          <Button
            leftSection={<IconRobot size={16} />}
            onClick={() => setTestAIModal(true)}
            variant="light"
            color="grape"
          >
            اختبار AI
          </Button>
        </Group>
      </Group>

      {/* Filters */}
      <Card withBorder mb="md">
        <Group justify="space-between">
          <SegmentedControl
            value={filters.timeRange}
            onChange={(value) => setFilters({ ...filters, timeRange: value as any })}
            data={[
              { label: 'اليوم', value: 'today' },
              { label: 'أسبوع', value: 'week' },
              { label: 'شهر', value: 'month' },
              { label: 'ربع سنوي', value: 'quarter' },
              { label: 'سنوي', value: 'year' },
            ]}
          />
          
          <Group>
            {filters.timeRange === 'custom' && (
              <>
                <DatePicker
                  placeholder="من تاريخ"
                  value={filters.startDate ? new Date(filters.startDate) : null}
                  onChange={(date) => setFilters({ ...filters, startDate: date?.toISOString() })}
                />
                <DatePicker
                  placeholder="إلى تاريخ"
                  value={filters.endDate ? new Date(filters.endDate) : null}
                  onChange={(date) => setFilters({ ...filters, endDate: date?.toISOString() })}
                />
              </>
            )}
            
            <Select
              placeholder="المنصة"
              data={[
                { label: 'Instagram', value: 'instagram' },
                { label: 'TikTok', value: 'tiktok' },
                { label: 'Twitter', value: 'twitter' },
                { label: 'YouTube', value: 'youtube' },
                { label: 'Facebook', value: 'facebook' },
                { label: 'AI Services', value: 'ai' },
              ]}
              value={filters.platform}
              onChange={(value) => setFilters({ ...filters, platform: value || undefined })}
              clearable
            />
          </Group>
        </Group>
      </Card>

      {/* Tabs */}
      <Tabs value={selectedTab} onChange={setSelectedTab} mb="md">
        <Tabs.List>
          <Tabs.Tab value="overview" leftSection={<IconDashboard size={16} />}>
            نظرة عامة
          </Tabs.Tab>
          <Tabs.Tab value="ai" leftSection={<IconRobot size={16} />}>
            الذكاء الاصطناعي
          </Tabs.Tab>
          <Tabs.Tab value="services" leftSection={<IconServer size={16} />}>
            الخدمات
          </Tabs.Tab>
          <Tabs.Tab value="users" leftSection={<IconUsers size={16} />}>
            المستخدمين
          </Tabs.Tab>
          <Tabs.Tab value="orders" leftSection={<IconShoppingCart size={16} />}>
            الطلبات
          </Tabs.Tab>
          <Tabs.Tab value="analytics" leftSection={<IconChartBar size={16} />}>
            التحليلات
          </Tabs.Tab>
          <Tabs.Tab value="monitoring" leftSection={<IconAlertCircle size={16} />}>
            المراقبة
          </Tabs.Tab>
        </Tabs.List>

        {/* Overview Tab */}
        <Tabs.Panel value="overview" pt="md">
          <SimpleGrid cols={{ base: 1, sm: 2, lg: 4 }} spacing="md" mb="md">
            {renderMetricCard(
              'إجمالي المستخدمين',
              dashboardData?.stats.totalUsers || 0,
              dashboardData?.stats.growthRate,
              <IconUsers size={20} />
            )}
            
            {renderMetricCard(
              'إجمالي الإيرادات',
              adminHelpers.formatCurrency(dashboardData?.stats.totalRevenue || 0),
              undefined,
              <IconTrendingUp size={20} />
            )}
            
            {renderMetricCard(
              'طلبات الذكاء الاصطناعي',
              dashboardData?.stats.totalAIRequests || 0,
              undefined,
              <IconRobot size={20} />
            )}
            
            {renderMetricCard(
              'طلبات توليد الفيديو',
              dashboardData?.stats.videoGenerationRequests || 0,
              undefined,
              <IconVideo size={20} />
            )}
          </SimpleGrid>

          {/* Revenue Chart */}
          <Card withBorder mb="md">
            <Group justify="space-between" mb="md">
              <Title order={4}>الإيرادات عبر الزمن</Title>
              <Badge variant="light" color="blue">
                {filters.timeRange}
              </Badge>
            </Group>
            <Box style={{ height: 300 }}>
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={revenueChartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke={theme.colors.dark[3]} />
                  <XAxis dataKey="period" stroke={theme.colors.gray[6]} />
                  <YAxis stroke={theme.colors.gray[6]} />
                  <Tooltip 
                    contentStyle={{ 
                      backgroundColor: theme.colors.dark[7],
                      borderColor: theme.colors.dark[4],
                      color: theme.white 
                    }}
                  />
                  <Legend />
                  <Line 
                    type="monotone" 
                    dataKey="revenue" 
                    stroke={theme.colors.blue[6]} 
                    name="الإيرادات الكلية"
                  />
                  <Line 
                    type="monotone" 
                    dataKey="aiRevenue" 
                    stroke={theme.colors.grape[6]} 
                    name="إيرادات الذكاء الاصطناعي"
                  />
                </LineChart>
              </ResponsiveContainer>
            </Box>
          </Card>

          {/* AI Usage Chart */}
          <Grid gutter="md">
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Card withBorder>
                <Title order={4} mb="md">استخدام الذكاء الاصطناعي</Title>
                <Box style={{ height: 300 }}>
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={aiUsageChartData}
                        cx="50%"
                        cy="50%"
                        labelLine={false}
                        label={({ category, percent }) => `${category}: ${(percent * 100).toFixed(0)}%`}
                        outerRadius={80}
                        fill="#8884d8"
                        dataKey="count"
                      >
                        {aiUsageChartData.map((entry, index) => (
                          <Cell key={`cell-${index}`} fill={[
                            theme.colors.blue[6],
                            theme.colors.grape[6],
                            theme.colors.pink[6],
                            theme.colors.cyan[6],
                            theme.colors.green[6],
                          ][index % 5]} />
                        ))}
                      </Pie>
                      <Tooltip />
                    </PieChart>
                  </ResponsiveContainer>
                </Box>
              </Card>
            </Grid.Col>

            <Grid.Col span={{ base: 12, md: 6 }}>
              <Card withBorder>
                <Title order={4} mb="md">أداء النظام</Title>
                <Box style={{ height: 300 }}>
                  <ResponsiveContainer width="100%" height="100%">
                    <RadarChart data={performanceChartData}>
                      <PolarGrid />
                      <PolarAngleAxis dataKey="metric" />
                      <PolarRadiusAxis />
                      <Radar
                        name="Cloudflare"
                        dataKey="value"
                        stroke={theme.colors.blue[6]}
                        fill={theme.colors.blue[6]}
                        fillOpacity={0.6}
                      />
                    </RadarChart>
                  </ResponsiveContainer>
                </Box>
              </Card>
            </Grid.Col>
          </Grid>
        </Tabs.Panel>

        {/* AI Tab */}
        <Tabs.Panel value="ai" pt="md">
          <Card withBorder mb="md">
            <Group justify="space-between" mb="md">
              <Title order={4}>مقاييس الذكاء الاصطناعي</Title>
              <Group>
                <Badge color="grape" variant="light">
                  {dashboardData?.aiMetrics.totalRequests || 0} طلب
                </Badge>
                <Badge color="blue" variant="light">
                  {dashboardData?.aiMetrics.totalCost || 0} دولار تكلفة
                </Badge>
              </Group>
            </Group>

            <SimpleGrid cols={{ base: 1, md: 2, lg: 4 }} spacing="md">
              {renderMetricCard(
                'طلبات Gemini',
                dashboardData?.aiMetrics.gemini.requests || 0,
                undefined,
                <IconRobot size={20} color={theme.colors.grape[6]} />
              )}
              
              {renderMetricCard(
                'طلبات OpenAI',
                dashboardData?.aiMetrics.openai.requests || 0,
                undefined,
                <IconRobot size={20} color={theme.colors.blue[6]} />
              )}
              
              {renderMetricCard(
                'طلبات توليد الفيديو',
                dashboardData?.aiMetrics.videoGeneration.totalRequests || 0,
                undefined,
                <IconVideo size={20} color={theme.colors.pink[6]} />
              )}
              
              {renderMetricCard(
                'معدل النجاح',
                `${((dashboardData?.aiMetrics.totalSuccessfulRequests || 0) / (dashboardData?.aiMetrics.totalRequests || 1) * 100).toFixed(1)}%`,
                undefined,
                <IconCircleCheck size={20} color={theme.colors.green[6]} />
              )}
            </SimpleGrid>
          </Card>

          {/* AI Models Status */}
          <Card withBorder>
            <Title order={4} mb="md">حالة نماذج الذكاء الاصطناعي</Title>
            <Table>
              <thead>
                <tr>
                  <th>النموذج</th>
                  <th>الحالة</th>
                  <th>طلبات اليوم</th>
                  <th>معدل النجاح</th>
                  <th>متوسط الوقت</th>
                  <th>التكلفة</th>
                  <th>الإجراءات</th>
                </tr>
              </thead>
              <tbody>
                {Object.entries(dashboardData?.aiMetrics || {}).map(([key, model]: [string, any]) => (
                  typeof model === 'object' && model.requests !== undefined && (
                    <tr key={key}>
                      <td>
                        <Group gap="xs">
                          <Badge variant="light" color="grape">
                            {adminHelpers.formatAIModelName(key)}
                          </Badge>
                        </Group>
                      </td>
                      <td>
                        <Badge 
                          color={model.successRate > 90 ? 'green' : model.successRate > 70 ? 'yellow' : 'red'}
                          variant="light"
                        >
                          {model.successRate > 90 ? 'ممتاز' : model.successRate > 70 ? 'جيد' : 'تحتاج تحسين'}
                        </Badge>
                      </td>
                      <td>{model.requests || 0}</td>
                      <td>
                        <Group>
                          <Text>{model.successRate?.toFixed(1) || 0}%</Text>
                          <Progress value={model.successRate || 0} size="sm" w={100} />
                        </Group>
                      </td>
                      <td>{model.averageTime?.toFixed(2) || 0} ثانية</td>
                      <td>{adminHelpers.formatCurrency(model.cost || 0, 'USD')}</td>
                      <td>
                        <Group gap="xs">
                          <ActionIcon variant="light" size="sm" color="blue">
                            <IconEye size={14} />
                          </ActionIcon>
                          <ActionIcon variant="light" size="sm" color="yellow">
                            <IconPlayerPlay size={14} />
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

        {/* Services Monitoring Tab */}
        <Tabs.Panel value="services" pt="md">
          <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="md">
            {services.map((service) => (
              <Card key={service.name} withBorder>
                <Group justify="space-between" mb="md">
                  <div>
                    <Text fw={500}>{service.name}</Text>
                    <Text size="sm" c="dimmed">{service.description}</Text>
                  </div>
                  <Badge 
                    color={getServiceStatusColor(service)}
                    variant="light"
                  >
                    {service.status || 'unknown'}
                  </Badge>
                </Group>
                
                <Stack gap="xs">
                  <Group justify="space-between">
                    <Text size="sm">الاستجابة</Text>
                    <Text fw={500}>{service.responseTime || 0}ms</Text>
                  </Group>
                  <Group justify="space-between">
                    <Text size="sm">التشغيل</Text>
                    <Text fw={500}>{service.uptime || 0}%</Text>
                  </Group>
                  <Group justify="space-between">
                    <Text size="sm">الأخطاء</Text>
                    <Badge color={service.errorRate > 5 ? 'red' : 'green'}>
                      {service.errorRate || 0}%
                    </Badge>
                  </Group>
                </Stack>

                <Button 
                  fullWidth 
                  mt="md" 
                  variant="light" 
                  leftSection={<IconRefresh size={16} />}
                  onClick={() => {
                    // Refresh service logic here
                  }}
                >
                  تحديث الخدمة
                </Button>
              </Card>
            ))}
          </SimpleGrid>
        </Tabs.Panel>

        {/* Monitoring Tab */}
        <Tabs.Panel value="monitoring" pt="md">
          <Grid gutter="md">
            {/* System Health */}
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Card withBorder>
                <Title order={4} mb="md">صحة النظام</Title>
                {systemHealth && (
                  <Stack gap="md">
                    <Group justify="space-between">
                      <Text>الحالة العامة</Text>
                      <Badge 
                        color={
                          systemHealth.status === 'healthy' ? 'green' :
                          systemHealth.status === 'degraded' ? 'yellow' : 'red'
                        }
                      >
                        {systemHealth.status === 'healthy' ? 'صحي' :
                         systemHealth.status === 'degraded' ? 'متدهور' : 'غير صحي'}
                      </Badge>
                    </Group>

                    <Stack gap="xs">
                      {Object.entries(systemHealth.components).map(([component, status]) => (
                        <Group key={component} justify="space-between">
                          <Text size="sm">{component}</Text>
                          <Badge 
                            color={status ? 'green' : 'red'}
                            variant="light"
                            leftSection={status ? 
                              <IconCircleCheck size={12} /> : 
                              <IconCircleX size={12} />
                            }
                          >
                            {status ? 'نشط' : 'غير نشط'}
                          </Badge>
                        </Group>
                      ))}
                    </Stack>
                  </Stack>
                )}
              </Card>
            </Grid.Col>

            {/* Recent Alerts */}
            <Grid.Col span={{ base: 12, md: 6 }}>
              <Card withBorder>
                <Title order={4} mb="md">التنبيهات الأخيرة</Title>
                <ScrollArea style={{ height: 400 }}>
                  <Timeline active={0}>
                    {dashboardData?.systemAlerts.slice(0, 10).map((alert) => (
                      <Timeline.Item 
                        key={alert.id}
                        title={alert.title}
                        bullet={
                          alert.severity === 'critical' ? 
                            <IconAlertCircle size={16} color={theme.colors.red[6]} /> :
                            <IconBellRinging size={16} color={theme.colors.yellow[6]} />
                        }
                      >
                        <Text size="sm" c="dimmed">{alert.message}</Text>
                        <Text size="xs" mt={4}>
                          {new Date(alert.timestamp).toLocaleString('ar-SA')}
                        </Text>
                        {alert.actionRequired && (
                          <Button size="xs" variant="light" color="red" mt="xs">
                            إجراء مطلوب
                          </Button>
                        )}
                      </Timeline.Item>
                    ))}
                  </Timeline>
                </ScrollArea>
              </Card>
            </Grid.Col>
          </Grid>
        </Tabs.Panel>
      </Tabs>

      {/* Export Modal */}
      <Modal 
        opened={exportModal} 
        onClose={() => setExportModal(false)}
        title="تصدير تقرير"
        size="lg"
      >
        <Stack gap="md">
          <Select
            label="نوع التقرير"
            data={[
              { label: 'الطلبات', value: 'orders' },
              { label: 'المستخدمين', value: 'users' },
              { label: 'الإيرادات', value: 'revenue' },
              { label: 'التحليلات', value: 'analytics' },
              { label: 'الذكاء الاصطناعي', value: 'ai' },
              { label: 'الكل', value: 'all' },
            ]}
            defaultValue="orders"
          />
          
          <MultiSelect
            label="التنسيقات"
            data={[
              { label: 'PDF', value: 'pdf' },
              { label: 'Excel', value: 'excel' },
              { label: 'CSV', value: 'csv' },
              { label: 'JSON', value: 'json' },
            ]}
            defaultValue={['pdf', 'excel']}
          />
          
          <Group>
            <Switch label="تضمين الرسوم البيانية" defaultChecked />
            <Switch label="تضمين التفاصيل" defaultChecked />
          </Group>
          
          <Select
            label="المنطقة الزمنية"
            data={[
              { label: 'السعودية (توقيت الرياض)', value: 'Asia/Riyadh' },
              { label: 'GMT', value: 'GMT' },
              { label: 'UTC', value: 'UTC' },
            ]}
            defaultValue="Asia/Riyadh"
          />
          
          <Button 
            fullWidth 
            leftSection={<IconDownload size={16} />}
            onClick={() => {
              notifications.show({
                title: 'جاري التصدير',
                message: 'سيتم إعداد التقرير وإرساله إليك',
                color: 'blue',
              });
              setExportModal(false);
            }}
          >
            تصدير الآن
          </Button>
        </Stack>
      </Modal>

      {/* Test AI Modal */}
      <Modal 
        opened={testAIModal} 
        onClose={() => setTestAIModal(false)}
        title="اختبار خدمة الذكاء الاصطناعي"
        size="lg"
      >
        <Stack gap="md">
          <Select
            label="نموذج الذكاء الاصطناعي"
            data={[
              { label: 'Google Gemini', value: 'gemini' },
              { label: 'OpenAI GPT', value: 'openai' },
              { label: 'Ollama LLM', value: 'ollama' },
              { label: 'Hugging Face', value: 'huggingface' },
              { label: 'Stability AI (صور)', value: 'stability' },
              { label: 'Luma AI (فيديو)', value: 'luma' },
              { label: 'Runway ML (فيديو)', value: 'runway' },
              { label: 'Pika Labs (فيديو)', value: 'pika' },
              { label: 'Gemini Veo (فيديو)', value: 'gemini-veo' },
            ]}
            value={selectedModel}
            onChange={(value) => value && setSelectedModel(value)}
          />
          
          <Textarea
            label="النص (اختياري)"
            placeholder="أدخل نص للاختبار..."
            value={testPrompt}
            onChange={(e) => setTestPrompt(e.target.value)}
            autosize
            minRows={3}
          />
          
          <Group justify="space-between">
            <Text size="sm" c="dimmed">
              التكلفة المتوقعة: {adminHelpers.estimateAICost(selectedModel, 1, 'text')} دولار
            </Text>
            
            <Button 
              loading={testAIMutation.isPending}
              leftSection={<IconPlayerPlay size={16} />}
              onClick={() => testAIMutation.mutate()}
            >
              بدء الاختبار
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Box>
  );
};

export default AdminDashboard;