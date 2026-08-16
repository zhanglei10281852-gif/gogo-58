# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

做审计演练时发现 audit.jsonl 的校验拦不住中间被动过的记录：我们把三条记录里的第二条换成另一个 store 同位置的记录，序号仍然是 1、2、3，每条记录自身的摘要也都算得上，结果 verify 照样说链有效并给出 head，完全看不出中间那条已经不是原来的记录了。只有直接删掉一行、把序号弄断号时才会被发现。请修复审计链校验，让中间记录被替换或重排（序号连续、单条自洽）时也判定为无效并指出断点，同时空日志仍视为有效、正常日志仍报告正确的 head、自身摘要被改的记录仍在原位置被发现，并保证全量测试通过。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-58
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-58.git
- parent SHA：47261516ab8d3a43514a11b19e7d5245e58c7fac

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-58.git bug-repro
cd bug-repro
git checkout --detach 47261516ab8d3a43514a11b19e7d5245e58c7fac
go test ./internal/store -run "^TestVerifyAuditRejectsSplicedChain$" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/store -run "^TestVerifyAuditRejectsSplicedChain$" -count=1 -v
=== RUN   TestVerifyAuditRejectsSplicedChain
    audit_link_regression_test.go:46: a spliced audit log verified as a valid chain: {EntryCount:3 Head:aa17927a2e890026b5985890f94b21a3768e79269a247f48222348987ef9db66 Valid:true BrokenAt:0 Reason:}
--- FAIL: TestVerifyAuditRejectsSplicedChain (0.07s)
FAIL
FAIL	CaveLoop/internal/store	0.078s
FAIL

```

stderr：

```text
warning: internal/store/audit_link_regression_test.go has type 100755, expected 100644
warning: internal/store/audit_link_regression_test.go has type 100755, expected 100644

```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/store -run "^TestVerifyAuditRejectsSplicedChain$" -count=1 -v
=== RUN   TestVerifyAuditRejectsSplicedChain
    audit_link_regression_test.go:46: a spliced audit log verified as a valid chain: {EntryCount:3 Head:aa17927a2e890026b5985890f94b21a3768e79269a247f48222348987ef9db66 Valid:true BrokenAt:0 Reason:}
--- FAIL: TestVerifyAuditRejectsSplicedChain (0.10s)
FAIL
FAIL	CaveLoop/internal/store	0.220s
FAIL

```

stderr：

```text
warning: internal/store/audit_link_regression_test.go has type 100755, expected 100644
warning: internal/store/audit_link_regression_test.go has type 100755, expected 100644

```

## 通过条件

把一条来自其它审计日志的记录拼进链中（seq 连续 1..3、每条自身摘要合法）时 VerifyAudit 返回 valid=false、brokenAt=2、head 为空；未被改动的三条链仍 valid=true、entryCount=3、head 非空；空日志仍 valid=true 且 head=GenesisHash；detail 被篡改的记录仍在 brokenAt=1 报“does not match”；删除一行导致序号断号的既有判定不回归；定向测试、全量 go test ./... -count=1 与 go build ./... && go vet ./... 全部通过；校准与远端复跑均在 golang:1.22 linux/amd64 单架构完成。
