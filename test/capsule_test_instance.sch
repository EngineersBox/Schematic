variable "testcapsule_clsid" {
  value = 85831
}

instance "capsule::config" "test_capsule" {
  inbuilt = true
  containerId = "test_capsule_id"
  config = {
    pidsMax = 20
    memMax = 4096
    netClsId = var.testcapsule_clsid
    terminateOnClose = true
  }
}